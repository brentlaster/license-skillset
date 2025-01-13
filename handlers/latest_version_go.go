package handlers

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"os"
)

func LatestVersionGo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Latest Go Version Called")
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://go.dev/dl/?mode=json&include=all", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Println("Something unrecoverable went wrong")
		return
	}

	// Parse the JSON response
	var versions []GoVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
	   return "", err
	}
	
	// The latest version is the first item in the array
	if len(versions) > 0 {
	   return versions[0].Version, nil
	}
	return "", fmt.Errorf("no versions found")
}
