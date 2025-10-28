package commands

import "strings"

type Kind int

const (
	KindRegistry Kind = iota
	KindGitHub
)

func getKind(arg string) Kind {
	var githubPrefixes = []string{
		"https://github.com/",
		"http://github.com/",
		"github.com/",
	}

	var prefix string
	for _, prefix = range githubPrefixes {
		if strings.HasPrefix(arg, prefix) {
			return KindGitHub
		}
	}

	return KindRegistry
}
