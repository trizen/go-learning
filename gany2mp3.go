/*
 * gany2mp3
 * Copyright (C) 2014 Daniel "Trizen" Șuteu
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

/*
   gany2mp3: convert any file to MP3 by running multiple instances
             of ffmpeg at the same time in a controllable way.

   Date: 12 April 2014 [17:10:49]
   License: GPLv3
   Website: http://github.com/trizen

   Compilation:
       go build gany2mp3.go

   Usage:
       ./gany2mp3 [options] [files]

   Options:
       -t int      : use this many threads (default: 2)
       -o str      : output directory (default: .)
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func executeCmd(channel chan string, name string, arg []string) {
	cmd := exec.Command(name, arg...)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stdout, err)
	}
	channel <- "."
}

func main() {
	const maxFormatLen = 1 + 4 // .webm
	const outputFormat = "mp3"

	const ffmpegCmd = "ffmpeg"

	ffmpegArg := []string{"-y", "-vn", "-ac", "2", "-ab", "192K", "-ar", "48000", "-f", outputFormat}

	var outputDir string
	var maxThreads int

	// Get the flags
	flag.StringVar(&outputDir, "o", ".", "put converted files in this directory")
	flag.IntVar(&maxThreads, "t", 2, "the number of threads used for conversion")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "error: no input file provided!")
		os.Exit(2)
	}

	// Get the input files
	inputFiles := flag.Args()

	// Create a channel and a counter
	channel := make(chan string)
	counter := 0

	for i, file := range inputFiles {

		outputFile := file

		// Remove the format suffix
		if index := strings.LastIndex(outputFile, "."); index != 1 && len(outputFile)-index <= maxFormatLen {
			outputFile = outputFile[0:index]
		}

		// Trim the parent directory
		if index := strings.LastIndex(outputFile, "/"); index != 1 {
			outputFile = outputFile[index+1:]
		}

		outputFile = outputDir + "/" + outputFile + "." + outputFormat
		fmt.Printf("[%2d] %s -> %s\n", i+1, file, outputFile)

		args := []string{"-i", file}
		args = append(args, ffmpegArg...)
		args = append(args, outputFile)

		// Execute the ffmpeg command (in a goroutine)
		go executeCmd(channel, ffmpegCmd, args)

		// Don't allow more than maxThreads commands
		// to run in parallel
		if counter += 1; counter >= maxThreads {
			fmt.Print("** Waiting...")
			fmt.Println(<-channel)
		}
	}

	fmt.Print("** Waiting...")

	for i := 1; i < maxThreads; i++ {

		// When more threads has been specified than input files
		if i > counter {
			break
		}

		fmt.Print(<-channel)
	}

	fmt.Println("All done!")
}
