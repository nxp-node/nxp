package api_github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/nxp-node/nxp/pkg/api"
)

var ErrMissingManifest = errors.New("missing manifest")
var ErrMissingMetadata = errors.New("missing metadata")

func QueryManifest(repo string) (mod *api.Manifest, err error) {
	url := "https://raw.githubusercontent.com/" + repo + "/master/package.json"
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}

	if resp.Status == "404" {
		return nil, ErrMissingManifest
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &mod); err != nil {
		return nil, err
	}

	return mod, nil
}

func QueryMetadata(repo string) (mod *Metadata, err error) {
	url := "https://api.github.com/repos/" + repo

	req, err := http.NewRequest("GET", url, http.NoBody)
	req.Header.Add("Accept", "application/vnd.github+json")

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}

	if resp.Status == "404" {
		return nil, ErrMissingMetadata
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &mod); err != nil {
		return nil, err
	}

	return mod, nil
}

func DownloadPackage(repo string, versionName string, version string) (targz *[]byte, suggestedFilename string, err error) {
	fn := fmt.Sprintf("%s-%s.tgz", versionName, version)
	url := fmt.Sprintf("https://github.com/%s/archive/%s.tar.gz", repo, version)

	resp, err := http.Get(url)
	if err != nil {
		// panic(fmt.Sprintf("Failed to download the module 'github.com/%s'", repo))
		return nil, fn, err
	} else {
		defer resp.Body.Close()
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fn, err
	}

	return &data, fn, nil
}
