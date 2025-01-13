package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GoVersion struct {
	Version string `json:"version"`
}

func LatestVersionGo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Latest Go Version Called")
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://go.dev/dl/?mode=json", nil)
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

	// we need to use tee here because io.Copy consumes the resp.Body stream and leads to EOF for decoding
	var bodyBuffer bytes.Buffer
	tee := io.TeeReader(resp.Body, &bodyBuffer)
	if _, err := io.Copy(w, tee); err != nil {
		fmt.Println("Something unrecoverable went wrong")
		return
	}

	// Parse the JSON response
	var versions []GoVersion
	if err := json.NewDecoder(&bodyBuffer).Decode(&versions); err != nil {
		fmt.Println(err)
		fmt.Println(resp.Body)
		fmt.Println("Could not decode response")
		return
	}

	// The latest version is the first item in the array
	if len(versions) > 0 {
		fmt.Println(versions[0].Version)
		return
	}
	fmt.Println("No versions found")

}
