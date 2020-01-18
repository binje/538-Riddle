package main

import "fmt"
import "math/rand"

var edges map[int][]int = make(map[int][]int)

func main() {
	addEdges()

	fmt.Println("Two ducks simulation expected hops:")
	fmt.Println(simulations(2))

	fmt.Println("Three ducks simulation expected hops:")
	//fmt.Println(simulations(3))

	fmt.Println("Four ducks simulation expected hops:")
	//fmt.Println(simulations(4))

	// I know that this would be easy using matrix operations but I really wated to do this the long way around
	solveAnalytical()

}

func solveAnalytical() {

	equations := make(map[point]equation)

	for i := 1; i <= 9; i++ {
		for j := i; j <= 9; j++ {
			for k := j; k <= 9; k++ {
				equations[point{i, j, k}] = createEquation(i, j, k)
			}
		}
	}

	for pnt, _ := range equations {
		// All points that are the same that are not the start should be removed
		if pnt.X == pnt.Y && pnt.Y == pnt.Z && pnt.X != 5 {
			delete(equations, pnt)
		}
		// Ducks always move odd<->even so these states are not possible
		if pnt.X%2 != pnt.Y%2 || pnt.X%2 != pnt.Z%2 {

			delete(equations, pnt)
		}
	}
	for p, _ := range equations {
		if (p == point{5, 5, 5}) {
			continue
		}
		removeEquation(equations, p)
	}
	for k, v := range equations {
		fmt.Println(k, v, v.Constant)
	}
	fmt.Println(equations[point{5, 5, 5}].Constant)
}

func createEquation(x, y, z int) equation {
	points := make([]point, 0)
	for _, rockX := range edges[x] {
		for _, rockY := range edges[y] {
			for _, rockZ := range edges[z] {
				// Always have X < Y < Z to remove redundancy
				if rockX > rockY {
					rockX, rockY = rockY, rockX
				}
				if rockY > rockZ {
					rockY, rockZ = rockZ, rockY
				}
				if rockX > rockY {
					rockX, rockY = rockY, rockX
				}
				points = append(points, point{rockX, rockY, rockZ})
			}
		}
	}
	denom := int64(len(points))
	phrases := make(map[point]fraction)
	constant := fraction{1, 1}
	for _, pnt := range points {
		if pnt.X != pnt.Y || pnt.Y != pnt.Z {
			f := phrases[pnt]
			f.plus(fraction{1, denom})
			phrases[pnt] = f
		}
	}
	return equation{phrases, &constant}
}

func removeEquation(eqs map[point]equation, p point) {
	if f, ok := eqs[p].Phrases[p]; ok {
		left := fraction{1, 1}
		left.minus(f)
		inverse := fraction{left.Denominator, left.Numerator}
		for pnt, frac := range eqs[p].Phrases {
			frac.mul(inverse)
			eqs[p].Phrases[pnt] = frac
		}
		eqs[p].Constant.mul(inverse)
		delete(eqs[p].Phrases, p)
	}

	for _, eq := range eqs {
		k, ok := eq.Phrases[p]
		if !ok {
			continue
		}
		for pnt, frac := range eqs[p].Phrases {
			frac.mul(k)
			f := eq.Phrases[pnt]
			f.plus(frac)
			eq.Phrases[pnt] = f
		}
		k.mul(*eqs[p].Constant)
		f := eq.Constant
		f.plus(k)
		eq.Constant = f
		delete(eq.Phrases, p)
	}
	delete(eqs, p)
}

type equation struct {
	//LeftHand point
	Phrases  map[point]fraction
	Constant *fraction
}

type point struct {
	X, Y, Z int
}

type fraction struct {
	Numerator, Denominator int64
}

func (f *fraction) mul(f2 fraction) {
	f.Numerator *= f2.Numerator
	f.Denominator *= f2.Denominator
	f.simplify()
}

func (f *fraction) plus(f2 fraction) {
	if f.Numerator == 0 && f.Denominator == 0 {
		f.Numerator = f2.Numerator
		f.Denominator = f2.Denominator
		return
	}
	f.Numerator = f.Numerator*f2.Denominator + f.Denominator*f2.Numerator
	f.Denominator = f.Denominator * f2.Denominator
	f.simplify()
}

func (f *fraction) minus(f2 fraction) {
	f.plus(fraction{-1 * f2.Numerator, f2.Denominator})
}

func (f *fraction) simplify() {
	gcd := GCD(f.Numerator, f.Denominator)
	f.Numerator /= gcd
	f.Denominator /= gcd
}

func GCD(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

func addEdges() {
	edges[1] = []int{2, 4}
	edges[2] = []int{1, 3, 5}
	edges[3] = []int{2, 6}
	edges[4] = []int{1, 5, 7}
	edges[5] = []int{2, 4, 6, 8}
	edges[6] = []int{3, 5, 9}
	edges[7] = []int{4, 8}
	edges[8] = []int{5, 7, 9}
	edges[9] = []int{6, 8}
}

func simulations(numDucks int) float64 {
	n := 1000000
	total := 0
	for i := 0; i < n; i++ {
		total += simulate(numDucks)
	}
	return float64(total) / float64(n)
}

func simulate(numDucks int) int {
	ducks := make([]int, numDucks)
	for i := 0; i < numDucks; i++ {
		ducks[i] = hop(5)
	}
	steps := 1
	for notAllSameRock(ducks) {
		for i := 0; i < numDucks; i++ {
			ducks[i] = hop(ducks[i])
		}
		steps++
	}
	return steps
}

func notAllSameRock(rocks []int) bool {
	for i := 1; i < len(rocks); i++ {
		if rocks[i] != rocks[i-1] {
			return true
		}
	}
	return false
}

func hop(rock int) int {
	return edges[rock][rand.Intn(len(edges[rock]))]
}
