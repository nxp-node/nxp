package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nxp-node/nxp/cmd/console"
	"github.com/nxp-node/nxp/pkg/api"

	github "github.com/nxp-node/nxp/pkg/api_github"
	registry "github.com/nxp-node/nxp/pkg/api_registry"
)

func Install(args []string) {
	var arg = args[0]
	var kind = getKind(arg)

	var name string
	var version string
	var dependencies map[string]string

	var err error
	var tardata *[]byte
	var tarname string

	if kind == KindGitHub {
		var mod *api.Manifest

		url := arg
		repo := strings.SplitN(url, ".com/", 2)[1]

		mod, err = github.QueryManifest(repo)

		if err != nil {
			console.Printnln(Prefix + "⚠  error: couldn't query the manifest for the specified package")
			console.Printnln(Prefix + "           " + err.Error())
			return
		}

		name = mod.Name
		version = mod.Version

		if mod.Dependencies != nil {
			dependencies = *mod.Dependencies
		}

		console.Printf(Prefix+"📩 downloading %s@%s", name, version)
		tardata, tarname, err = github.DownloadPackage(repo, name, version)
	} else {
		var mod *registry.Module

		name = arg
		mod, err = registry.QueryPackage(name)

		if err != nil {
			console.Printnln(Prefix + "⚠  error: couldn't query the specified package")
			console.Printnln(Prefix + "           " + err.Error())
			return
		}

		version = mod.DistTags["latest"]
		manifest := mod.Versions[version]

		if manifest.Dependencies != nil {
			dependencies = *manifest.Dependencies
		}

		console.Printf(Prefix+"📩 downloading %s@%s", name, version)
		tardata, tarname, err = registry.DownloadPackage(name, manifest.Name, version)
	}

	if err != nil {
		console.Println(Prefix + "⚠  error: couldn't download the specified package")
		console.Println(Prefix + "           " + err.Error())
		return
	}

	dir := "./nxp_modules"
	os.Mkdir(dir, 0700)

	unxtracted := dir + "/package"
	tarPath := dir + "/" + tarname
	destination := dir + "/" + arg

	if _, err = os.Stat(unxtracted); err == nil {
		console.Printnln(Prefix + "⚠  error: decompression folder ('package') already exists")
		console.Printnln(Prefix + "           (to continue, rename or delete it)")
		return
	}

	if overwriteCheck(tarname, Option{'s', "skip", "skipping"}) == OC_OPT(1) {
		console.Print(Prefix + "🛌 skipping package installation")
		return
	}

	var destOpt = overwriteCheck(destination, Option{'u', "update", "updating"}, Option{'s', "skip", "skipping"})
	if destOpt == OC_OPT(0) {
		var val api.Manifest

		data, err := os.ReadFile(destination + "/manifest.json")
		if err != nil {
			console.Println(Prefix + "⚠  error: couldn't read the local package's manifest")
			console.Println(Prefix + "⚠          " + err.Error())
		}

		err = json.Unmarshal(data, &val)
		if err != nil {
			console.Println(Prefix + "⚠  error: couldn't parse the local package's manifest")
			console.Println(Prefix + "⚠          " + err.Error())
		}

		outdated := isNewer(val.Version, version)
		if !outdated {
			console.Print(Prefix + "📅 package is not outdated, skipping") // or 🎿?
			return
		}
	} else if destOpt == OC_OPT(1) {
		console.Print(Prefix + "🛌 skipping package installation")
		return
	}

	tgOpt := overwriteCheck(tarPath, Option{'s', "skip", "skipping"})
	if tgOpt == OC_OPT(0) {
		console.Print(Prefix + "🛌 skipping package installation")
		return
	}

	tgPrinted := tgOpt == OC_SAFE
	if tgPrinted {
		console.Print(Prefix + "📝 writing .tar.gz")
	}

	if err = os.WriteFile(tarPath, *tardata, 0700); err != nil {
		if tgPrinted {
			console.Print("\n")
		}

		console.Println(Prefix + "⚠  error: couldn't create the specified package's tar.gz")
		console.Println(Prefix + "           " + err.Error())
		return
	}

	if len(dependencies) > 0 && tgPrinted {
		console.Print("\n")
	}

	for dependency, version := range dependencies {
		console.Printf(Prefix+"🖋️  installing dependency: %s[#568856]%s[/#568856]", dependency, version)
		Install([]string{
			dependency,
		})
	}

	console.Print(Prefix + "🤐 extracting")

	var cmd = exec.Cmd{
		Path: "tar",
		Args: []string{"tar", "-xf", tarname},
		Dir:  dir,
	}

	lp, _ := exec.LookPath("tar")
	if lp != "" {
		cmd.Path = lp
	}

	if err = cmd.Run(); err != nil {
		console.Print(Prefix + "⚠  error: couldn't extract the specified package's tar.gz")
		console.Print(Prefix + "           " + err.Error())
		return
	}

	if _, err = os.Stat(destination); err == nil {
		console.Print(Prefix + "🚫 deleting old version of the package")
		os.Rename(destination, fmt.Sprintf("%s/nxp-lost-%d---%s", os.TempDir(), time.Now().UnixMilli(), arg))
	}

	console.Print(Prefix + "📦 renaming decompressed folder ('package') to the package name")
	os.Rename(unxtracted, destination)

	console.Print(Prefix + "🚫 deleting .tar.gz")
	os.Remove(tarPath)

	console.Print("")
}

func isNewer(current string, new string) bool {
	localVersion := strings.Split(current, ".")
	remoteVersion := strings.Split(new, ".")

	for i, num := range remoteVersion {
		if len(localVersion)-1 < i {
			return true
		}

		if localVersion[i] < num {
			return true
		}
	}

	return false
}
