package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"./build"
)

const (
	configsDir  = "configs"
	templateDir = "build/templates"
	outputDir   = "node"
)

func main() {
	if len(os.Args) < 2 {
		var coins []string
		filepath.Walk(filepath.Join(configsDir, "coins"), func(path string, info os.FileInfo, err error) error {
			n := strings.TrimSuffix(info.Name(), ".json")
			if n != info.Name() {
				coins = append(coins, n)
			}
			return nil
		})
		fmt.Fprintf(os.Stderr, "Usage: %s coin <url>\nCoin build-in are:\n%v\n", filepath.Base(os.Args[0]), coins)
		os.Exit(0)
	}

	coin := os.Args[1]
        var url string
        if len(os.Args) > 2 {
          url = os.Args[2]
        }

	config, err := build.LoadConfig(configsDir, coin, url )

	if err == nil {
		err = build.GeneratePackageDefinitions(config, templateDir, outputDir, coin)
	} else {
               fmt.Printf("%s\n", err)
               os.Exit(0)
        }
        err = os.RemoveAll(filepath.Join(outputDir, "backend/scripts"))
        if err != nil {
        	panic(err)
        }
        err = os.RemoveAll(filepath.Join(outputDir, "backend/config"))
        if err != nil {
                panic(err)
        }

        if config.Backend.ExecScript != "" {
                err = os.Remove(filepath.Join(outputDir, "backend/" + coin + ".conf"))
        	if err != nil {
                	panic(err)
        	}
        }



	fmt.Printf("Package files for %v generated to %v\n", coin, outputDir)
}
