package main

import "fmt"
import "net/http"
import "bufio"

//import "strings"
import "github.com/willf/bitset"

func main() {
	fmt.Println("vim-go")
	url := "https://norvig.com/ngrams/enable1.txt"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	allWords := make([]string, 0)
	for scanner.Scan() {
		allWords = append(allWords, scanner.Text())
	}

	sets := make(map[string]*letters)
	for _, word := range allWords {
		if len(word) < 4 {
			continue
		}
		b := getBitSet(word)
		count := b.Count()
		if count < 8 && !b.Test(uint('s'-'a')) {
			// use string as key because bitSet cannot be used
			s := b.DumpAsBits()
			if _, ok := sets[s]; !ok {
				sets[s] = &letters{b, 0}
			}
			// points for this word
			p := len(word)
			// 4 letter words are worth one point
			if p == 4 {
				p = 1
			}
			// Panagram
			if count == 7 {
				p += 7
			}
			sets[s].Points += p
		}
	}
	fmt.Println(len(sets))
	pangrams := make([]letters, 0)
	nonPangrams := make([]letters, 0)
	for _, v := range sets {
		if v.B.Count() == 7 {
			pangrams = append(pangrams, *v)
		} else {
			nonPangrams = append(nonPangrams, *v)
		}
	}
	fmt.Println(len(pangrams))

	max := 0
	var maxB *bitset.BitSet
	var maxL uint
	for _, pangram := range pangrams {
		for i := uint(0); i < 26; i++ {
			p := pangram.Points
			if !pangram.B.Test(i) {
				continue
			}
			for _, word := range nonPangrams {
				if word.B.Test(i) && pangram.B.IsSuperSet(word.B) {
					p += word.Points
				}
			}
			if p > max {
				max = p
				maxB = pangram.B
				maxL = i
			}
		}
	}
	fmt.Println(max)
	fmt.Println(maxB)
	fmt.Println(maxL)
	for _, word := range allWords {
		if len(word) < 4 {
			continue
		}
		b := getBitSet(word)
		if b.Test(maxL) && maxB.IsSuperSet(b) {
			fmt.Println(word)
		}
	}
}

type letters struct {
	B      *bitset.BitSet
	Points int
}

func getBitSet(word string) *bitset.BitSet {
	b := bitset.New(26)
	for _, c := range word {
		b.Set(uint(c - 'a'))
	}
	return b
}
