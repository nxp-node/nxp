package api

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// package.json
type Manifest struct {
	Name            string             `json:"name"`
	Version         string             `json:"version"`
	Description     string             `json:"description"`
	Keywords        *[]string          `json:"keywords"`
	Author          *any               `json:"author"`
	Dependencies    *map[string]string `json:"dependencies"`
	DevDependencies *map[string]string `json:"devDependencies"`
	Bin             *map[string]string `json:"bin"`
	Main            *string            `json:"main"`
	Engines         *map[string]string `json:"engines"`
}

type Module struct {
	ID             string              `json:"id"`
	Revision       string              `json:"_rev"`
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	DistTags       map[string]string   `json:"dist-tags"`
	Versions       map[string]Manifest `json:"versions"`
	Maintaners     []User              `json:"maintaners"`
	Time           map[string]string   `json:"time"`
	Author         User                `json:"author"`
	Users          map[string]bool     `json:"users"`
	Readme         string              `json:"readme"`
	ReadmeFilename string              `json:"readmeFilename"`
	Keywords       []string            `json:"keywords"`
}
