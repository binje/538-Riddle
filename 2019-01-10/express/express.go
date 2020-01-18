package main

import "fmt"

func main() {
	fmt.Println("vim-go")
	for denom := 1; denom < 10000; denom++ {
		num := (denom / 2020) + 1
		if denom <= 2019*num {
			continue
		}
		fmt.Prinln(num)
		fmt.Println(denom)
		break
	}
}
