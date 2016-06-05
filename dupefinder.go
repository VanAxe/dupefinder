package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app				= kingpin.New("dupefinder", "Find duplicate files with checksums.")
	verbose		= app.Flag("verbose", "Enable verbose mode.").Short('v').Bool()
	debug			= app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	rec				= app.Flag("recursive", "Search recursively.").Short('r').Bool()
	output		= app.Flag("output", "File to which the report will be written.").Short('o').Default("./dupefinder.log").String()
	target		= app.Arg("target", "Where to look for duplicate files.").Default(".").String()
)

func main() {
	app.Version("0.0.1")

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug == true {
		conf := `== DEBUG CONF ==
Verbose:		%t
Debug:			%t
Recursive:		%t
Output:			%s
Target:			%s

`
		fmt.Printf(conf, *verbose, *debug, *rec, *output, *target)
	}
	find_dupes()
}

func find_dupes() {
	if *rec == true {
		err := filepath.Walk(*target, visit)
		fmt.Printf("filepath.Walk() returned %v\n", err)
	} else {
		files, err := ioutil.ReadDir(*target)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
}

// func isDir(pth string) (bool, error) {
// 	fi, err := os.Stat(pth)
// 	if err != nil {
// 		return false, err
// 	}

// 	return fi.Mode.IsDir(), nil
// }

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s\n", path)
  return nil
}
