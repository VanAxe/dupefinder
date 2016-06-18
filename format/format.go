package format

import (
	"fmt"
)

func Pp_csv(input_map map[string][]string) string {
	pretty_csv := ""
	for entry, files := range input_map {
		string_of_files := ""
		for _, file := range files {
			string_of_files = string_of_files + ", " + file
		}
		pretty_csv = pretty_csv + entry + string_of_files + "\n"
	}
	return pretty_csv
}

func Pp_human(input_map map[string][]string) string {
	pretty_csv := ""
	// for entry, files := range input_map {
	// 	string_of_files := ""
	// 	for _, file := range files {
	// 		string_of_files = string_of_files + ", " + file
	// 	}
	// 	pretty_csv = pretty_csv + entry + string_of_files + "\n"
	// 	// fmt.Printf("%s, %s\n ", entry, string_of_files)
	// }
	return pretty_csv
}

func Print_debug(verbose bool, debug *bool, rec *bool, output string, abs_path string) {
	if *debug == true {
		fmt.Println("== DEBUG CONF ==")
		fmt.Println("Verbose:\t", verbose)
		fmt.Println("Debug:\t\t", *debug)
		fmt.Println("Recursive:\t", *rec)
		fmt.Println("Output:\t\t", output)
		fmt.Println("Target:\t\t", abs_path, "\n")
	}
}
