package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
        "log"
        "io/ioutil"
	"github.com/trezor/blockbook/build/tools"
)

const (
	configsDir  = "configs"
	templateDir = "build/templates"
	outputDir   = "pkg-build"
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
		fmt.Fprintf(os.Stderr, "Usage: %s coin\nCoin is one of:\n%v\n", filepath.Base(os.Args[0]), coins)
		os.Exit(1)
	}

	coin := os.Args[1]

        if coin == "list" {
        	files, err := ioutil.ReadDir("configs/coins")
        	if err != nil {
          		log.Fatal(err)
        	}

    		for _, file := range files {
        	fmt.Println(strings.ReplaceAll(file.Name(), ".json", ""))
    		}
          os.Exit(0)
        }


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
//        err = os.RemoveAll(filepath.Join(outputDir, "backend/debian"))
//        if err != nil {
//                panic(err)
//        }
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
