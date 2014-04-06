/*
   fsfn - Find Similar File Names
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
   Find similar file names in a given directory tree

   Date: 06 April 2014
   Website: http://github.com/trizen

   Compilation:
       go build fsfn.go

   Usage:
       ./fsfn [dir]
*/

package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

// Compare two strings and return true
// if they look pretty much the same.
func alikeStr(x, y string) bool {

	if x == y {
		return true
	}

	len1, len2 := len(x), len(y)

	if len1 > len2 {
		y, x, len2, len1 = x, y, len1, len2
	}

	min := int(math.Ceil(float64(len2) / 2))
	if min > len1 {
		return false
	}

	diff := len1 - min

	for i := 0; i <= diff; i++ {
		for j := i; j <= diff; j++ {
			if strings.Index(y, x[i:i+(min+j-i)]) != -1 {
				return true
			}
		}
	}

	return false
}

// Traverse a given directory path and
// return a list of pairs [path, modName]
func traverse(dir string) (files [][2]string) {
	dirData, err := ioutil.ReadDir(dir)

	if err != nil {
		return
	}

	for _, file := range dirData {
		if fname := file.Name(); fname[0] != '.' {

			if file.IsDir() {
				files = append(files, traverse(dir+"/"+fname)...)
			} else {
				name := strings.ToLower(fname)
				idx := strings.LastIndex(name, ".")

				if idx != -1 && len(name)-idx <= 5 {
					files = append(files, [2]string{dir + "/" + fname, name[0:idx]})
				} else {
					files = append(files, [2]string{dir + "/" + fname, name})
				}
			}
		}
	}

	return
}

// This is the main function
func main() {

	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "[dir]")
		return
	}

	data := traverse(dir)

	files := map[string][]string{}
	for i := 0; i < len(data)-1; i++ {
		for j := i + 1; j < len(data); j++ {
			if alikeStr(data[i][1], data[j][1]) {
				files[data[i][0]] = append(files[data[i][0]], data[j][0])
				dataBuff := data[0:j]
				data = append(dataBuff, data[j+1:]...)
				j--
			}
		}
	}

	for key, val := range files {
		newArr := append(val, key)
		sort.Strings(newArr)
		for _, path := range newArr {
			fmt.Println(path)
		}
		fmt.Println(strings.Repeat("-", 80))
	}
}
