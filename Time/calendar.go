package main

import (
	"fmt"
	"time"
)

func main() {
	mons := map[int]int{
		1:  31,
		2:  28,
		3:  31,
		4:  30,
		5:  31,
		6:  30,
		7:  31,
		8:  31,
		9:  30,
		10: 31,
		11: 30,
		12: 31,
	}

	t := time.Now()
	cur_mon := int(t.Month())
	cur_year := t.Year()

	cur_mon_name := fmt.Sprintf("%s", t.Month())
	fmt.Printf("%*s\n%s\n", (20+5+len(cur_mon_name))/2, fmt.Sprintf("%s %d", cur_mon_name, cur_year), "Su Mo Tu We Th Fr Sa")

	if cur_year%400 == 0 || cur_year%4 == 0 && cur_year%100 != 0 {
		mons[2] = 29
	}

	cur_year -= 1
	st := int(1 + cur_year*365 + cur_year/4 - cur_year/100 + cur_year/400)

	for i := 1; i < int(t.Month()); i++ {
		st += mons[i]
	}

	for i := 1; i <= st%7; i++ {
		fmt.Printf("   ")
	}
	cur_year += 1

	for i := 1; i <= mons[cur_mon]; i++ {
		fmt.Printf("%2d ", i)

		if (st+i)%7 == 0 && i != mons[cur_mon] {
			fmt.Println()
		}
	}

	fmt.Println("\n")
}
