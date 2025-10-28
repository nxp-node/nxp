package commands

import (
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
	var err error
	var tardata *[]byte
	var tarname string

	prefixpre := "[#1e5688]nxp[/#1e5688] [#336997]"
	prefixsuf := "[/#336997] "

	prefix := prefixpre + "│" + prefixsuf

	if kind == KindGitHub {
		var mod *api.Manifest

		url := arg
		repo := strings.SplitN(url, ".com/", 2)[1]

		mod, err = github.QueryManifest(repo)

		if err != nil {
			console.Println(prefix + "⚠  error: couldn't query the manifest for the specified package")
			console.Println(prefix + "           " + err.Error())
			return
		}

		name = mod.Name
		version := mod.Version

		console.Printf(prefix+"📩 downloading %s@%s", name, version)
		tardata, tarname, err = github.DownloadPackage(repo, name, version)
	} else {
		var mod *registry.Module

		name = arg
		mod, err = registry.QueryPackage(name)

		if err != nil {
			console.Println(prefix + "⚠  error: couldn't query the specified package")
			console.Println(prefix + "           " + err.Error())
			return
		}

		latest := mod.DistTags["latest"]
		version := mod.Versions[latest]

		console.Printf(prefix+"📩 downloading %s@%s", name, latest)
		tardata, tarname, err = registry.DownloadPackage(name, version.Name, latest)
	}

	if err != nil {
		console.Println(prefix + "⚠  error: couldn't download the specified package")
		console.Println(prefix + "           " + err.Error())
		return
	}

	dir := "./nxp_modules"
	os.Mkdir(dir, 0700)

	unxtracted := dir + "/package"
	tarPath := dir + "/" + tarname
	destination := dir + "/" + arg

	if _, err = os.Stat(unxtracted); err == nil {
		console.Println(prefix + "⚠  error: decompression folder ('package') already exists")
		console.Println(prefix + "           (to continue, rename or delete it)")
		return
	}

	overwriteCheck(tarname)
	overwriteCheck(destination)

	if !overwriteCheck(tarPath) {
		console.Print(prefix + "📝 writing .tar.gz")
	}

	if err = os.WriteFile(tarPath, *tardata, 0700); err != nil {
		console.Println(prefix + "⚠  error: couldn't create the specified package's tar.gz")
		console.Println(prefix + "           " + err.Error())
		return
	}

	console.Print(prefix + "🤐 extracting")

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
		console.Print(prefix + "⚠  error: couldn't extract the specified package's tar.gz")
		console.Print(prefix + "           " + err.Error())
		return
	}

	if _, err = os.Stat(destination); err == nil {
		console.Print(prefix + "🚫 deleting old version of the package")
		os.Rename(destination, fmt.Sprintf("%s/nxp-lost-%d---%s", os.TempDir(), time.Now().UnixMilli(), arg))
	}

	console.Print(prefix + "📦 renaming decompressed folder ('package') to the package name")
	os.Rename(unxtracted, destination)

	console.Print(prefix + "🚫 deleting .tar.gz")
	os.Remove(tarPath)

	console.Print("")
}
