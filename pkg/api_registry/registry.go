package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Olafcio1/nxp/pkg/api"
)

var ErrNotFound = errors.New("package not found")
var ErrUnrecognized = errors.New("unrecognized error")

func QueryPackage(name string) (mod *api.Module, err error) {
	url := "https://registry.npmjs.org/" + name
	resp, err := http.Get(url)

	if err != nil {
		panic(fmt.Sprintf("Failed to query the module '%s'", name))
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	str := string(data)
	if strings.HasPrefix(str, "{\"error\":") {
		var errorText = str[10 : len(str)-2]

		if errorText == "Not found" {
			return nil, ErrNotFound
		} else {
			return nil, ErrUnrecognized
		}
	}

	if err := json.Unmarshal(data, &mod); err != nil {
		return nil, err
	}

	return mod, nil
}

func DownloadPackage(name string, versionName string, version string) (targz *[]byte, suggestedFilename string, err error) {
	fn := fmt.Sprintf("%s-%s.tgz", versionName, version)
	url := fmt.Sprintf("https://registry.npmjs.org/%s/-/%s", name, fn)

	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("Failed to download the module '%s'", name))
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fn, err
	}

	return &data, fn, nil
}
