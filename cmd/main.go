package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/mod/modfile"
)

var (
	// flag for repository A
	repoA = flag.String("repoA", "", "repository A")
	// flag for repository B
	repoB = flag.String("repoB", "", "repository B")
)

func init() {
	flag.Parse()
}

func main() {
	if *repoA == "" || *repoB == "" {
		flag.PrintDefaults()
		return
	}

	repoAModFile, err := getModFile(*repoA)
	if err != nil {
		fmt.Printf("Error getting repo A mod file at URL %s: %v\n", *repoA, err)
		os.Exit(1)
	}

	repoBModFile, err := getModFile(*repoB)
	if err != nil {
		fmt.Printf("Error getting repo B mod file at URL %s: %v\n", *repoB, err)
		os.Exit(1)
	}

	fmt.Println("Comparing mod files...")

	modfileA, err := modfile.Parse("go.mod", repoAModFile, nil)
	if err != nil {
		fmt.Printf("Error parsing repo A mod file: %v\n", err)
		os.Exit(1)
	}

	modfileB, err := modfile.Parse("go.mod", repoBModFile, nil)
	if err != nil {
		fmt.Printf("Error parsing repo B mod file: %v\n", err)
		os.Exit(1)
	}

	modfileAPathsAndVersions := getPathsAndVersions(modfileA)
	modfileBPathsAndVersions := getPathsAndVersions(modfileB)

	for path, version := range modfileAPathsAndVersions {
		if modBVersion, ok := modfileBPathsAndVersions[path]; ok {
			if modBVersion != version {
				fmt.Printf("Your module requires %s@%s, the comparison requires %s\n", path, version, modBVersion)
			}
		}
	}

	os.Exit(0)
}

func getModFile(repo string) ([]byte, error) {
	resp, err := http.Get(repo)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func getPathsAndVersions(modfile *modfile.File) map[string]string {
	pathsAndVersions := make(map[string]string, len(modfile.Require))

	for _, require := range modfile.Require {
		pathsAndVersions[require.Mod.Path] = require.Mod.Version
	}

	return pathsAndVersions
}
