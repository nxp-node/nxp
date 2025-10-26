package cmd

type Subcommand struct {
	Usage         string
	ArgumentCount Range
	Description   string
	Function      func([]string)
}
