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
	Author          any                `json:"author"`
	Dependencies    *map[string]string `json:"dependencies"`
	DevDependencies *map[string]string `json:"devDependencies"`
	Bin             *map[string]string `json:"bin"`
	Main            *string            `json:"main"`
	Engines         *map[string]string `json:"engines"`
}
