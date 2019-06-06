package main

/*

Author: Trizen
License: GPLv3
Date: 02 November 2012
Website: http://trizen.googlecode.com

This program generates a valid pipe-menu XML for the Openbox Window Manager.

This is a partial translation of the Perl script: obmenu-generator.
See: http://code.google.com/p/trizen/downloads/detail?name=obmenu-generator

Usage: ./gobmenugen

Write in your ~/.config/openbox/menu.xml file, the following text:

<?xml version="1.0" encoding="utf-8"?>
<openbox_menu xmlns="http://openbox.org/"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://openbox.org/">
    <menu id="root-menu" label="obmenu-generator" execute="PATH_TO_THIS_COMPILED_PROGRAM" />
</openbox_menu>

After that, execute 'openbox --reconfigure'.
For more settings, see the configuration file: ~/.config/gobmenugen/config

*/

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func write_config_file(config_file string) {
	stat, err := os.Stat(config_file)
	config_content := `# SCHEMA supports the following keys: item, cat, bgcat, encat, exit, raw, sep
#
# Posible values for each of this types are:
# For 'item'  : name,command
# For 'sep'   : A string representing the LABEL for the separator.
# For 'cat'   : name,catname - Any of the possible categories.
# For 'exit'  : name
# For 'raw'   : A hardcoded XML line in the Openbox's menu.xml file format
#    Example  : raw:<menu icon="" id="client-list-combined-menu" />
# For 'bgcat' : label,id        (begin a category)
# For 'encat' : [empty string]  (end of a category)

[CONFIG]
terminal:xterm

[SCHEMA]
item:File Manager,pcmanfm
item:Terminal,xterm
item:Editor,geany
sep:Applications
cat:Accessories,Utility
cat:Development,Development
cat:Education,Education
cat:Games,Games
cat:Graphics,Graphics
cat:Multimedia,AudioVideo
cat:Network,Network
cat:Office,Office
cat:Settings,Settings
cat:System,System
sep:
exit:Exit
`

	if err != nil || stat.Size() == 0 {
		out_fh, err := os.OpenFile(config_file, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			log.Fatal(err)
		}

		defer out_fh.Close()

		_, err = out_fh.WriteString(config_content)

		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

func print_header() {
	fmt.Println("<openbox_pipe_menu>")
}

func print_item(name, command string) {
	fmt.Printf("    <item label=\"%s\"><action name=\"Execute\"><execute>%s</execute></action></item>\n", name, command)
}

func print_category(cat_id, label string) {
	fmt.Printf("  <menu id=\"%s\" label=\"%s\">\n", cat_id, label)
}

func print_end_category() {
	fmt.Println("  </menu>")
}

func print_exit(label string) {
	fmt.Printf("  <item label=\"%s\"><action name=\"Exit\" /></item>\n", label)
}

func print_footer() {
	fmt.Println("</openbox_pipe_menu>")
}

func main() {

	username, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	home := username.HomeDir
	config_file := filepath.Join(home, ".config", "gobmenugen", "config")
	stat, err := os.Stat(config_file)

	if err != nil || stat == nil {
		dir := filepath.Dir(config_file)
		err := os.MkdirAll(dir, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	write_config_file(config_file)

	fh, err := os.Open(config_file)

	if err != nil {
		log.Fatal(err)
	}

	wantedCats := make(map[string]int)
	bf := bufio.NewReader(fh)

	var (
		config, schema bool
	)

	config_data := make(map[string]string)
	schema_data := make(map[string][]string)
	var schema_keys []string

	for {
		line, isPrefix, err := bf.ReadLine()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if string(line) == "" {
			continue
		}

		switch line[0] {
		case '#':
			continue
		}

		if string(line) == "[SCHEMA]" {
			schema = true
			config = false
			continue
		} else if string(line) == "[CONFIG]" {
			schema = false
			config = true
			continue
		}

		if isPrefix {
			log.Fatal("Error: Unexpected long line reading", fh.Name())
		}

		parts := bytes.SplitN(line, []byte{':'}, 2)

		key := string(parts[0])

		if schema == true {
			schema_data[key] = append(schema_data[key], string(parts[1]))
			schema_keys = append(schema_keys, key)
			if key == "cat" {
				values := bytes.SplitN(parts[1], []byte{','}, 2)
				wantedCats[string(values[1])] = 1
			}
		} else if config == true {
			config_data[key] = string(parts[1])
		}
	}

	apps_dir := "/usr/share/applications"
	dir, err := ioutil.ReadDir(apps_dir)

	if err != nil {
		log.Fatal(err)
	}

	type Item struct {
		name, exec, icon, terminal string
	}

	apps := make(map[string][]Item)

	var f os.FileInfo
	for _, f = range dir {

		if f.Name()[len(f.Name())-8:] != ".desktop" {
			continue
		}

		file, err := ioutil.ReadFile(filepath.Join(apps_dir, f.Name()))

		if err != nil {
			continue
		}

		lines := bytes.Split(file, []byte{'\n'})

		var (
			name, exec, icon, terminal string
			categories                 []byte
			noDisplay                  bool
		)

		foundDE := false

		for _, line := range lines {
			line = bytes.TrimSpace(line)

			if string(line) == "" || string(line[0]) == "#" {
				continue
			}

			if foundDE == false {
				if string(line) == "[Desktop Entry]" {
					foundDE = true
				}
				continue
			}

			if string(line[0]) == "[" {
				break
			}

			parts := bytes.SplitN(line, []byte{'='}, 2)

			switch string(parts[0]) {
			case "Exec":
				exec = string(parts[1])
			case "Name":
				name = string(parts[1])
			case "Icon":
				icon = string(parts[1])
			case "Categories":
				categories = parts[1]
			case "Terminal":
				terminal = string(parts[1])
			}

			if string(parts[0]) == "NoDisplay" {
				val := string(parts[1])
				if val == "true" || val == "1" {
					noDisplay = true
					break
				}
			}
		}

		if noDisplay == true {
			continue
		}

		cats := bytes.Split(categories, []byte{';'})

		for _, catName := range cats {
			cat := string(catName)

			if cat == "" {
				continue
			}

			if wantedCats[cat] != 1 {
				continue
			}

			apps[cat] = append(apps[cat], Item{
				name,
				exec,
				icon,
				terminal,
			})
		}
	}

	tracing := make(map[string]int)

	print_header()
	for _, key := range schema_keys {

		value := schema_data[key][tracing[key]]

		switch key {

		case "sep":
			if value == "" {
				fmt.Println("  <separator/>")
			} else {
				fmt.Printf("  <separator label=\"%s\"/>\n", value)
			}

		case "cat":
			parts := strings.SplitN(value, ",", 2)
			items := apps[parts[1]]

			if len(items) == 0 {
				tracing[key] += 1
				continue
			}

			print_category(parts[1], parts[0])

			for _, hash_ref := range items {
				exec := hash_ref.exec

				if strings.Contains(exec, "%") == true {
					exec = strings.SplitN(exec, " %", 2)[0]
				}

				if strings.Contains(hash_ref.name, "&") == true {
					hash_ref.name = strings.Replace(hash_ref.name, "&", "&amp;", -1)
				}

				if strings.ToLower(hash_ref.terminal) == "true" || hash_ref.terminal == "1" {
					exec = fmt.Sprintf("%s -e '%s'", config_data["terminal"], exec)
				}

				print_item(hash_ref.name, exec)
			}

			print_end_category()

		case "item":
			values := strings.SplitN(value, ",", 2)
			print_item(values[0], values[1])

		case "raw":
			fmt.Println(value)

		case "exit":
			print_exit(value)

		case "bgcat":
			values := strings.SplitN(value, ",", 2)
			print_category(values[1], values[0])

		case "encat":
			print_end_category()
		}

		tracing[key] += 1
	}
	print_footer()
}
