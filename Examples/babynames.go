package main

/*

Author: Trizen
License: GPLv3
Date: 26 November 2012
Website: http://trizen.googlecode.com

See also: https://code.google.com/intl/ro/edu/languages/google-python-class/set-up.html

*/

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	summaryfile := flag.Bool("summaryfile", false, "Create a new .summary file")
	flag.Parse()

	get_year_re, err := regexp.Compile(`\bPopularity in (\d+)\b`)
	if err != nil {
		log.Fatal(err)
	}

	get_name_re, err := regexp.Compile(`<tr align="right"><td>(\d+)</td><td>(\S+)</td><td>(\S+)</td>`)
	if err != nil {
		log.Fatal(err)
	}

	var args = flag.Args()
	if len(args) == 0 {
		fmt.Println("usage: [--summaryfile] file [file ...]")
		os.Exit(1)
	}

	for _, file := range args {
		var fh, err = os.Open(file)

		if err != nil {
			log.Println(err)
			continue
		}

		bf := bufio.NewReader(fh)

		var year string
		var m = map[string]string{}

		for {
			line, err := bf.ReadString('\n')

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
			}

			if line == "" {
				continue
			}

			match_year := get_year_re.FindStringSubmatch(line)
			if match_year != nil {
				year = match_year[1]
			}

			if year == "" {
				continue
			}

			match_name := get_name_re.FindStringSubmatch(line)
			if match_name != nil {
				pos := match_name[1]
				names := match_name[2:]

				for _, name := range names {

					if m[name] != "" {
						old_pos, err := strconv.Atoi(m[name])
						if err != nil {
							log.Fatal(err)
						}
						new_pos, err := strconv.Atoi(pos)
						if err != nil {
							log.Fatal(err)
						}

						if old_pos > new_pos {
							m[name] = pos
						}
					} else {
						m[name] = pos
					}
				}

			} else {
				continue
			}
		}

		var values [][]string
		for key, value := range m {
			values = append(values, []string{key, value})
		}

		// Sort by name
		for index := range values {
			value := values[index]
			i := index - 1
			for {
				if i < 0 {
					break
				}
				if value[0] < values[i][0] {
					values[i+1] = values[i]
					values[i] = value
					i -= 1
				} else {
					break
				}
			}
		}

		var summary_fh *os.File
		if *summaryfile == true {

			// Create .summary file
			summary_fh, err = os.Create(fmt.Sprintf("%s.summary", file))
			if err != nil {
				log.Fatal(err)
			}

			// Write year
			_, err := summary_fh.WriteString(fmt.Sprintf("%s\n", year))
			if err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Println(year)
		}

		for _, pairs := range values {
			line := fmt.Sprintf("%s %s\n", pairs[0], pairs[1])

			if *summaryfile == true {
				_, err := summary_fh.WriteString(line)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Print(line)
			}
		}
	}
}
