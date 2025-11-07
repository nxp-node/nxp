package commands

import (
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
			console.Println(Prefix + "⚠  warning: couldn't fetch the last downloads number")
			// console.Println("                " + err.Error())
			// return

			var num = uint(8)
			downloads = &num
		}

		latest := mod.DistTags["latest"]
		manifest = mod.Versions[latest]

		stat = fmt.Sprintf("📈 %d", downloads)
	}

	if err != nil {
		console.Println(Prefix + "⚠  error: couldn't download the specified package")
		console.Println(Prefix + "           " + err.Error())
		return
	}

	console.Println(PrefixPre + "├──────────────┤" + PrefixSuf)
	console.Println(PrefixPre + "│ package view │" + PrefixSuf)
	console.Println(PrefixPre + "├──────────────┤" + PrefixSuf)
	console.Fprintln(
		Prefix+"%s@%s | %s",
		manifest.Name, manifest.Version, stat,
	)
	console.Fprintln(
		Prefix+"%s",
		manifest.Description,
	)

	console.Fprintln(Prefix)

	if manifest.Keywords != nil {
		console.Fprintln(
			Prefix+"keywords: %s",
			strings.Join(*manifest.Keywords, ", "),
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

	useGrad, _ := gradient.NewGradient("#1eb0ffff", "#2496d3ff")
	devGrad, _ := gradient.NewGradient("#1eff8bff", "#24d38aff")

	console.Println(Prefix + "dependencies:")
	console.PrintEntries(
		slices.Concat(
			lo.MapToSlice(*manifest.Dependencies, func(key string, value string) []string {
				key = useGrad.Apply(key) + RESET
				return []string{key, value}
			}),
			lo.MapToSlice(*manifest.DevDependencies, func(key string, value string) []string {
				key = devGrad.Apply(key) + RESET
				return []string{key, value}
			}),
		),
		Prefix,
	)

	console.Fprintln(Prefix)

	if author := manifest.Author; author != nil {
		if val, ok := author.(string); ok {
			console.Fprintln(Prefix+"author: %s", val)
		} else if val, ok := author.(map[string]any); ok {
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
