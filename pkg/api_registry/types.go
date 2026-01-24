package api_registry

import "github.com/nxp-node/nxp/pkg/api"

type Module struct {
	ID          string                  `json:"id"`
	Revision    string                  `json:"_rev"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	DistTags    map[string]string       `json:"dist-tags"`
	Versions    map[string]api.Manifest `json:"versions"`
	Maintaners  []api.User              `json:"maintaners"`

	//Time           map[string]string       `json:"time"`     also sometimes messes up, unused so commenting it out

	Author         any             `json:"author"`
	Users          map[string]bool `json:"users"`
	Readme         string          `json:"readme"`
	ReadmeFilename string          `json:"readmeFilename"`
	Keywords       []string        `json:"keywords"`
}
