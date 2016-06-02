package main

import (
	"fmt"
	// "os"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	app	:= kingpin.New("dupefinder", "Find duplicate files with checksums.")

	config := configureApp(app)

	if *target != "" {
		fmt.Printf("Such wows.")
	} else {
		fmt.Printf("Such nopes.")
	}

}

func configureApp(app *kingpin.Application) {
		verbose		:= app.Flag("verbose", "Enable verbose mode.").Short('v').Bool()
		debug			:= app.Flag("debug", "Enable debug mode.").Short('d').Bool()
		rec				:= app.Flag("recursive", "Search recursively.").Short('r').Bool()
		output		:= app.Flag("output", "File to which the report will be written.").Short('o').Default("./dupefinder.log").String()
		target		:= app.Arg("target", "Where to look for duplicate files.").Default(".").String()
		app.Version("0.0.1")

		type Config struct {
			verbose *bool
			debug 	*bool
			rec			*bool
			output	*string
			target	*string
		}

		config := Config{verbose, debug, rec, output, target}

		return config
}
