package main

import "fmt"

var state map[Pennies]State = make(map[Pennies]State)

func main() {
	state[Pennies{0, 0}] = Lose

	for i := 1; i <= 30; i++ {
		for j := 0; j <= i; j++ {
			p := Pennies{i, j}
			state[p] = p.state()
		}
	}

	losers := make(map[int]struct{})
	for p, s := range state {
		sum := p.A + p.B
		if sum < 20 || sum > 30 {
			continue
		}
		if s == Lose {
			losers[sum] = struct{}{}
		}
	}

	for p, s := range state {
		sum := p.A + p.B
		if sum == 20 {
fmt.Println(p,s)
		}
	}

	for i:=20; i <=30; i++ {
		if _,ok := losers[i]; !ok {
			fmt.Println(i)
		}
	}
}

func (p *Pennies) state() State {
	for _, m := range p.nextMoves() {
		s, ok := state[m]
		if !ok {
			fmt.Println("Couldn't find")
			fmt.Println(p, m)
		}
		if s == Lose {
			return Win
		}
	}
	return Lose
}

func (p *Pennies) nextMoves() (moves []Pennies) {
	// remove only from A
	for i := 0; i < p.A; i++ {
		if i >= p.B {
			moves = append(moves, Pennies{i, p.B})
		} else {
			moves = append(moves, Pennies{p.B, i})
		}
	}
	// remove only from B
	for i := 0; i < p.B; i++ {
		moves = append(moves, Pennies{p.A, i})
	}
	// rmove from both piles
	for i := 1; i <= p.B; i++ {
		moves = append(moves, Pennies{p.A - i, p.B - i})
	}
	return
}

type Pennies struct {
	A, B int
}

type State int

const (
	_ State = iota
	Win
	Lose
)
