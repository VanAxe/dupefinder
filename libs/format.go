package dupefinder.format

func Pp_csv(input_map map[string][]string) string {
	pretty_csv := ""
	for entry, files := range input_map {
		string_of_files := ""
		for _, file := range files {
			string_of_files = string_of_files + ", " + file
		}
		pretty_csv = pretty_csv + entry + string_of_files + "\n"
		// fmt.Printf("%s, %s\n ", entry, string_of_files)
	}
	return pretty_csv
}

func Pp_human(input_map map[string][]string) string {
	// pretty_csv := ""
	// for entry, files := range input_map {
	// 	string_of_files := ""
	// 	for _, file := range files {
	// 		string_of_files = string_of_files + ", " + file
	// 	}
	// 	pretty_csv = pretty_csv + entry + string_of_files + "\n"
	// 	// fmt.Printf("%s, %s\n ", entry, string_of_files)
	// }
	// return pretty_csv
}

func Print_debug(conf string, *verbose *bool, *debug *bool, *rec *bool, output string, abs_path string, debug bool) {
	if *debug == true {
		conf := `== DEBUG CONF ==
		Verbose:		%t
		Debug:			%t
		Recursive:		%t
		Output:			%s
		Target:			%s

		`
		fmt.Printf(conf, *verbose, *debug, *rec, output, abs_path)
	}
}
