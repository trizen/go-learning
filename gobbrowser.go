package main

/*
# Author: Trizen
# License: GPLv3
# Date: 30 December 2012
# http://trizen.googlecode.com

# Recursively browse filesystem through openbox3 pipe menus

#---------------------------------------------------
# Add the filename of the built program into menu.xml as a menu id pipe:
#   <menu id="FILE_BROWSER" label="Disk" execute="/path/to/binary/file"/>
#---------------------------------------------------
*/

import (
	"bufio"
	"os"
	"os/user"
	"sort"
	"strings"
)

const (
	cmd = "pcmanfm" // open files with this application
)

func esc_str(s string) string {
	return strings.Replace(strings.Replace(s, "&", "&amp;", -1), "\"", "&quote;", -1)
}

func quot_escape(s string) string {
	return strings.Replace(s, "&quot;", "\\&quot;", -1)
}

func make_path(path, name string) string {
	return quot_escape(path + "/" + name)
}

func mk_dir_elem(path, name string) string {
	dir := make_path(path, name)
	return "<menu id=\"" + path + "/" + name + "\" label=\"" + name + "\" execute=\"" + os.Args[0] + " &quot;" + dir + "&quot;\"/>"
}

func mk_file_elem(path, name string) string {
	dir := make_path(path, name)
	return "<item label=\"" + name + "\"><action name=\"Execute\"><execute>" + cmd + " &quot;" + dir + "&quot;</execute></action></item>"
}

func main() {
	stdout := os.Stdout

	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		user, _ := user.Current()
		dir = user.HomeDir
	}

	var dirs, files []string
	dir_h, _ := os.Open(dir)
	defer dir_h.Close()

	dir_content, _ := dir_h.Readdir(-1)
	for _, file := range dir_content {

		if file.Name()[0] == '.' {
			continue
		}

		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else {
			files = append(files, file.Name())
		}
	}

	sort.Strings(files)
	sort.Strings(dirs)

	buf := bufio.NewWriter(stdout)

	escaped_dir := esc_str(dir)
	buf.WriteString("<openbox_pipe_menu>" + "<item label=\"Browse here...\"><action name=\"Execute\"><execute>" + cmd + " &quot;" + quot_escape(escaped_dir) + "&quot;</execute></action></item><separator/>")

	for _, file := range files {
		buf.WriteString(mk_file_elem(escaped_dir, esc_str(file)))
	}

	for _, dir := range dirs {
		buf.WriteString(mk_dir_elem(escaped_dir, esc_str(dir)))
	}
	buf.WriteString("</openbox_pipe_menu>")
	buf.Flush()
}
