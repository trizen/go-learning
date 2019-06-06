package main

/*
# Author: Daniel "Trizen" Șuteu
# License: GPLv3
# Date: 06 April 2014
# http://github.com/trizen

# Inspired from: Cosmos.A.Space.Time.Odyssey.S01E01
#                            by Neil deGrasse Tyson
*/

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// A <Pair> structure
type Pair struct {
	era  float64
	name string
}

// Method Pair.output(float64)
// Write to the standard output the content of a <Pair>
func (p *Pair) output(num float64) {
	fmt.Printf("\n=> In the Cosmic Year, that happened about %.2f %s ago!\n\n", num/p.era, p.name)
}

// Type <Pairs> -- an array of <Pair>s
type Pairs []Pair

// Method Pairs.push(float64, string)
// Append a new <Pair>, based on the last <Pair> of the array
func (p *Pairs) push(num float64, name string) {
	*p = append(*p, Pair{(*p)[len(*p)-1].era / num, name})
}

func main() {

	// Define the Cosmic Year
	rand.Seed(time.Now().UnixNano())
	cosmic_year := Pairs{Pair{(13.798 + []float64{+0.037, -0.037}[rand.Intn(2)]) * math.Pow(10.0, 9.0), "years"}}

	cosmic_year.push(12.0, "months")
	cosmic_year.push(30.4368499, "days")
	cosmic_year.push(24.0, "hours")
	cosmic_year.push(60.0, "minutes")
	cosmic_year.push(60.0, "seconds")
	cosmic_year.push(1000.0, "miliseconds")

	fmt.Println(`
This program will scale down the age of the universe to a normal year.

You can insert any number you want and the program will map it
into this cosmic year, giving you a feeling about how long ago
that happenned, compared to the age of the universe.
`)

	for {
		var year float64
		fmt.Print("> How long ago? (any number, in years): ")

		// Take input
		_, err := fmt.Scanf("%f", &year)
		if err != nil {
			continue
		}

		// Map the value inside the Cosmic Year
		for _, pair := range cosmic_year {
			if year >= pair.era {
				pair.output(year)
				break
			}
		}
	}
}
