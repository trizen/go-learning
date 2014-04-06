/*
   openbox-pipefs - Recursively browse filesystem through openbox3 pipe menus
   Copyright (C) 2014 Daniel "Trizen" Șuteu

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
   Recursively browse filesystem through openbox3 pipe menus

   Date: 06 April 2014
   Website: http://github.com/trizen

   Compilation:
       go build openbox-pipefs.go

    Usage:
    ---------------------------------------------------
    Add a new entry in your menu.xml:
        <menu id="openbox-pipefs" label="Disk" execute="openbox-pipefs ."/>
    ---------------------------------------------------

    ---------------------------------------------------
    If you are using the 'obmenu-generator' program, add in schema.pl:
        {pipe => ["openbox-pipefs .", "Disk", "drive-harddisk"]},
    ---------------------------------------------------

    See also: http://trizenx.blogspot.ro/2012/12/obbrowser.html
*/

package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
)

const CMD = "pcmanfm" // Command to lauch files with
const QUOT = "&#34;"  // XML escaped quote

type PathContent struct {
	dirs  []string
	files []string
}

func readdir(dir string) (content PathContent) {
	dirData, err := ioutil.ReadDir(dir)

	if err != nil {
		return
	}

	for _, file := range dirData {
		if fname := file.Name(); fname[0] != '.' {
			if file.IsDir() {
				content.dirs = append(content.dirs, fname)
			} else {
				content.files = append(content.files, fname)
			}
		}
	}

	return
}

type Buffer string

func (p *Buffer) Write(b []byte) (int, error) {
	*p += Buffer(b)
	return len(b), nil
}

func xmlEscape(str string) string {
	var buf Buffer
	xml.EscapeText(&buf, []byte(str))
	return string(buf)
}

func escapeQuot(str string) string {
	if strings.Contains(str, QUOT) {
		return strings.Replace(str, QUOT, "\\"+QUOT, -1)
	}
	return str
}

func main() {

	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		os.Stderr.Write([]byte("usage: " + os.Args[0] + " [dir]\n"))
		os.Exit(1)
	}

	content := readdir(dir)
	thisDir := xmlEscape(dir)
	qEscapedDir := escapeQuot(thisDir)
	escapedProgramName := xmlEscape(os.Args[0])

	os.Stdout.Write([]byte(
		"<openbox_pipe_menu><item label=\"Browse here...\">" +
			"<action name=\"Execute\"><execute>" + CMD + " " + QUOT + qEscapedDir + QUOT + "</execute></action>" +
			"</item><separator/>"))

	for _, name := range content.files {
		escapedName := xmlEscape(name)

		os.Stdout.Write([]byte(
			"<item label=\"" + escapedName + "\">" +
				"<action name=\"Execute\">" +
				"<execute>" + CMD + " " + QUOT + qEscapedDir + "/" + escapeQuot(escapedName) + QUOT + "</execute>" +
				"</action>" +
				"</item>"))
	}

	for _, name := range content.dirs {
		escapedName := xmlEscape(name)

		os.Stdout.Write([]byte(
			"<menu id=\"" + thisDir + "/" + escapedName + "\"" +
				" label=\"" + escapedName + "\"" +
				" execute=\"" + escapedProgramName + " " + QUOT + qEscapedDir + "/" + escapeQuot(escapedName) + QUOT +
				"\"/>"))
	}

	os.Stdout.Write([]byte("</openbox_pipe_menu>\n"))
}
