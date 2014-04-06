package main

/*
# CNP Info

# Author: Trizen
# License: GPLv3
# Date: 29 October 2012
# http://trizen.googlecode.com

# See:
# http://ro.wikipedia.org/wiki/Cod_numeric_personal
*/

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	var cnp = flag.Arg(0)

	if cnp == "" {
		fmt.Println("usage: ./cnp CNP")
		return
	}

	if len(cnp) != 13 {
		log.Fatal("Invalid CNP!")
	}

	var cnp_l = make([]int, len(cnp))

	for i := 0; i < len(cnp); i++ {
		num, err := strconv.Atoi(string(cnp[i]))
		if err != nil {
			log.Fatal(err)
		}
		cnp_l[i] = num
	}

	var magic = []int{2, 7, 9, 1, 4, 6, 3, 5, 8, 2, 7, 9}

	var jud = map[string]string{
		"01": "Alba",
		"02": "Arad",
		"03": "Argeș",
		"04": "Bacău",
		"05": "Bihor",
		"06": "Bistrița-Năsăud",
		"07": "Botoșani",
		"08": "Brașov",
		"09": "Brăila",
		"10": "Buzău",
		"11": "Caraș-Severin",
		"12": "Cluj",
		"13": "Constanța",
		"14": "Covasna",
		"15": "Dâmbovița",
		"16": "Dolj",
		"17": "Galați",
		"18": "Gorj",
		"19": "Harghita",
		"20": "Hunedoara",
		"21": "Ialomița",
		"22": "Iași",
		"23": "Ilfov",
		"24": "Maramureș",
		"25": "Mehedinți",
		"26": "Mureș",
		"27": "Neamț",
		"28": "Olt",
		"29": "Prahova",
		"30": "Satu Mare",
		"31": "Sălaj",
		"32": "Sibiu",
		"33": "Suceava",
		"34": "Teleorman",
		"35": "Timiș",
		"36": "Tulcea",
		"37": "Vaslui",
		"38": "Vâlcea",
		"39": "Vrancea",
		"40": "București",
		"41": "București S.1",
		"42": "București S.2",
		"43": "București S.3",
		"44": "București S.4",
		"45": "București S.5",
		"46": "București S.6",
		"51": "Călărași",
		"52": "Giurgiu",
	}

	type cnp_era struct {
		era int
		cet string
	}

	var year_era = map[string]cnp_era{
		"1": {1900, ""},
		"2": {1900, ""},
		"3": {1800, ""},
		"4": {1800, ""},
		"5": {2000, ""},
		"6": {2000, ""},
		"7": {
			0,
			"Străin rezident în România",
		},
		"8": {
			0,
			"Străin rezident în România",
		},
		"9": {
			0,
			"Persoană străină",
		},
	}

	var mon = []string{
		"Ianuarie",
		"Februarie",
		"Martie",
		"Aprilie",
		"Mai",
		"Iunie",
		"Iulie",
		"August",
		"Septembrie",
		"Octombrie",
		"Noiembrie",
		"Decembrie",
	}

	var days = []int{
		31,
		29,
		31,
		30,
		31,
		30,
		31,
		31,
		30,
		31,
		30,
		31,
	}

	var mon_map = map[string]int{}

	for i, name := range mon {
		mon_map[name] = days[i]
	}

	var sum = 0
	for i, _ := range magic {
		sum += magic[i] * cnp_l[i]
	}

	var cc = sum % 11
	if cc == 10 {
		cc = 1
	}

	if cc != cnp_l[12] {
		log.Fatal("Cifra de control e incorectă!\n")
	}

	yea_num := cnp_l[1]*10 + cnp_l[2]
	mon_num := cnp_l[3]*10 + cnp_l[4]
	day_num := cnp_l[5]*10 + cnp_l[6]
	jud_num := fmt.Sprintf("%d%d", cnp_l[7], cnp_l[8])

	var month_name = mon[mon_num-1]
	var jud_name = jud[jud_num]

	if jud_name == "" {
		log.Fatal("Codul județului e invalid!")
	}

	if day_num < 1 || day_num > mon_map[month_name] {
		log.Fatal("Ziua de naștere e invalidă!")
	}

	var hash_ref = year_era[strconv.Itoa(cnp_l[0])]
	var era = hash_ref.era

	cur_year, _, cur_day := time.Now().Date()

	t := time.Now()
	cur_mon, _ := strconv.Atoi(fmt.Sprintf("%d", t.Month()))

	var nationality = "Română"
	if era == 0 {
		if yea_num < cur_year-2000 {
			era = 2000
		} else {
			era = 1900
		}
		nationality = hash_ref.cet
	}

	birth_year := era + yea_num

	if mon_num < 1 || mon_num > 12 {
		log.Fatal("Luna de naștere e invalidă!\n")
	}

	if mon_num == 2 && day_num == 29 {
		if !(birth_year%400 == 0 || birth_year%4 == 0 && birth_year%100 != 0) {
			log.Fatal(fmt.Sprintf("Anul %d nu a fost un an bisect!\n", birth_year))
		}
	}

	var gender string
	if cnp_l[0] == 9 {
		gender = "Necunoscut"
	} else if cnp_l[0]%2 == 0 {
		gender = "Feminin"
	} else {
		gender = "Masculin"
	}

	var age = cur_year - birth_year
	if cur_mon < mon_num || (mon_num == cur_mon && day_num < cur_day) {
		age -= 1
	}

	fmt.Printf(
		`Data Nașterii:  %s
Cetațenie:      %s
Sexul:          %s
Vârsta:         %d
Județul:        %s
`,
		fmt.Sprintf("%d %s %d", day_num, month_name, birth_year), nationality, gender, age, jud_name)
}
