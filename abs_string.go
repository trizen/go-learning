package main

/*

# Author: Trizen
# License: GPLv3
# Date: 04 November 2012
# http://trizen.googlecode.com

# Expand a string to its absolute values

*/

import (
	"fmt"
	"strings"
)

func absolute_values(s string) []string {
	var abs []string
	var root []string

	parts := strings.Split(s, ",")

	append_item := func(chunk string) {
		abs = append(abs, strings.Join([]string{strings.Join(root, ""), chunk}, ""))
		return
	}

	for _, name := range parts {
		if strings.Contains(name, "{") == true {
			parts := strings.SplitAfter(name, "{")

			for _, chunk := range parts {
				if string(chunk[len(chunk)-1]) == "{" {
					root = append(root, chunk[0:len(chunk)-1])
				} else {
					append_item(chunk)
				}
			}

		} else if strings.Contains(name, "}") == true {
			parts := strings.Split(name, "}")

			for _, chunk := range parts {
				if len(chunk) == 0 {
					if len(root) != 0 {
						root = root[0 : len(root)-1]
					}
					continue
				}

				append_item(chunk)

			}
		} else {
			append_item(name)
		}
	}

	return abs
}

func main() {
	groups := []string{
		"perl-{gnome2-wnck,gtk2-{imageview,unique},x11-protocol,image-exiftool}",
		"perl-{proc-{simple,processtable},net-{dbus,dropbox-api},goo-canvas}",
		"perl-{sort-naturally,json,json-xs,xml-simple,www-mechanize,locale-gettext}",
		"perl-{file-{which,basedir,copy-recursive},pathtools,path-class},mplayer",
		"perl-{script-{test,meta}},flash-player",
	}
	for _, s := range groups {
		fmt.Printf("%+v\n", absolute_values(s))
	}
}
