package commands

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/iskaa02/qalam/gradient"
	"github.com/samber/lo"

	"github.com/nxp-node/nxp/cmd/console"
	"github.com/nxp-node/nxp/pkg/api"

	github "github.com/nxp-node/nxp/pkg/api_github"
	registry "github.com/nxp-node/nxp/pkg/api_registry"
)

func binCallback(value string, index int) string {
	slash := strings.LastIndex(value, "/") + 1
	filename := value[slash:]

	getter, _ := iter.Pull(strings.SplitSeq(filename, "."))

	name, _ := getter()
	name, _ = strings.CutPrefix(name, "node-core-")

	return name
}

func View(args []string) {
	var arg = args[0]
	var kind = getKind(arg)

	var stat string
	var name string
	var err error

	var manifest api.Manifest
	var isPackage = true

	if kind == KindGitHub {
		var mod *api.Manifest
		var info *github.Metadata

		url := arg
		repo := strings.SplitN(url, ".com/", 2)[1]

		mod, err = github.QueryManifest(repo)
		if err != nil {
			console.Println(Prefix + "⚠  error: couldn't query the manifest for the specified package repository")
			console.Println(Prefix + "           " + err.Error())
			return
		}

		info, err = github.QueryMetadata(repo)
		if err != nil {
			console.Println(Prefix + "⚠  error: couldn't query the metadata for the specified repository")
			console.Println(Prefix + "           " + err.Error())
			return
		}

		manifest = *mod
		stat = fmt.Sprintf("⭐ %d", info.StargazersCount)
	} else {
		var mod *registry.Module
		var downloads *uint

		name = arg
		mod, err = registry.QueryPackage(name)

		if err != nil {
			console.Println(Prefix + "⚠  error: couldn't query the specified package")
			console.Println(Prefix + "           " + err.Error())
			return
		}

		downloads, err = registry.GetLastMonth(name)

		if err != nil {
			if !errors.Is(err, registry.TypeError{}) {
				console.Println(Prefix + "⚠  warning: couldn't fetch the last downloads number")
			}
		} else {
			stat = fmt.Sprintf("📈 %d", downloads)
		}

		latest := mod.DistTags["latest"]
		manifest = mod.Versions[latest]

		if manifest.Name == "" {
			manifest = api.Manifest{
				Name:            mod.Name,
				Version:         latest,
				Description:     mod.Description,
				Keywords:        &mod.Keywords,
				Author:          mod.Author,
				Dependencies:    nil,
				DevDependencies: nil,
				Bin:             nil,
				Main:            nil,
			}

			isPackage = false
		}
	}

	if isPackage {
		console.Println(PrefixPre + "├──────────────┤" + PrefixSuf)
		console.Println(PrefixPre + "│ package view │" + PrefixSuf)
		console.Println(PrefixPre + "├──────────────┤" + PrefixSuf)
	} else {
		console.Println(PrefixPre + "├─────────────┤" + PrefixSuf)
		console.Println(PrefixPre + "│ module view │" + PrefixSuf)
		console.Println(PrefixPre + "├─────────────┤" + PrefixSuf)
	}

	console.Print(Prefix)
	fmt.Print(manifest.Name)
	if manifest.Version != "" {
		fmt.Print("@")
		fmt.Print(manifest.Version)
	}
	if stat != "" {
		fmt.Print(" | ")
		fmt.Print(stat)
	}

	console.Println("")

	if manifest.Description != "" {
		console.Fprintln(
			Prefix+"%s",
			manifest.Description,
		)
	}

	if manifest.Bin != nil {
		if binDict, ok := manifest.Bin.(map[string]string); ok {
			bin := lo.Map(
				lo.Values(binDict),
				binCallback,
			)

			slices.Sort(bin)
			slices.Reverse(bin)

			console.Fprintln(
				Prefix+"bin: %s",
				strings.Join(bin, ", "),
			)
		} else if binStr, ok := manifest.Bin.(string); ok {
			console.Fprintln(
				Prefix+"bin: ",
				binCallback(binStr, 0),
			)
		}
	}

	const RESET string = "\x1b[0m"

	if manifest.Dependencies != nil || manifest.DevDependencies != nil {
		useGrad, _ := gradient.NewGradient("#1eb0ffff", "#2496d3ff")
		devGrad, _ := gradient.NewGradient("#1eff8bff", "#24d38aff")

		if manifest.Dependencies != nil && len(*manifest.Dependencies) > 0 {
			console.Fprintln(Prefix)
			console.Println(Prefix + formatEach(useGrad.Apply("dependencies:"), "\x1b[1m", ""))
			console.PrintEntries(
				lo.MapToSlice(*manifest.Dependencies, func(key string, value string) []string {
					key = useGrad.Apply(key) + RESET
					return []string{key, value}
				}),
				Prefix,
			)
		}

		if manifest.DevDependencies != nil && len(*manifest.DevDependencies) > 0 {
			console.Fprintln(Prefix)
			console.Println(Prefix + formatEach(devGrad.Apply("dev dependencies:"), "\x1b[1m", ""))
			console.PrintEntries(
				lo.MapToSlice(*manifest.DevDependencies, func(key string, value string) []string {
					key = devGrad.Apply(key) + RESET
					return []string{key, value}
				}),
				Prefix,
			)
		}
	}

	printedKApref := false

	if manifest.Keywords != nil && len(*manifest.Keywords) > 0 {
		printedKApref = true

		console.Fprintln(Prefix)
		console.Fprintln(
			Prefix+"keywords: %s",
			strings.Join(*manifest.Keywords, ", "),
		)
	}

	if author := manifest.Author; author != nil && author != "" {
		if val, ok := author.(string); ok {
			if !printedKApref {
				console.Fprintln(Prefix)
			}

			console.Fprintln(Prefix+"author: %s", val)
		} else if val, ok := author.(map[string]any); ok {
			if !printedKApref {
				console.Fprintln(Prefix)
			}

			name, okN := val["name"].(string)
			email, okE := val["email"].(string)

			if okN || okE {
				out := ""

				if okN {
					out += " " + name
				}

				if okE {
					out += " <" + email + ">"
				}

				console.Fprintln(Prefix+"author:%s", out)
			}
		}
	}
}

func formatEach(text string, start string, end string) string {
	code := ""
	inCode := false

	out := ""
	for _, ch := range text {
		if inCode {
			if ch == 'm' {
				code += string(ch)
				out += code

				code = ""
				inCode = false
			} else {
				code += string(ch)
			}
		} else if ch == '\x1b' {
			inCode = true
			code = string(ch)
			out += start
		} else {
			out += string(ch) + end
		}
	}

	return out
}
