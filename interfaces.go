package main

import "fmt"

type animal interface {
	move() string
}

type dog struct{}
type fish struct{}
type bird struct{}

func (dog) move() string {
	return "I am a dog and walk"
}
func (fish) move() string {
	return "I am a fish and swim"
}
func (bird) move() string {
	return "I am a bird and fly"
}

func moveAnimal(a animal) {
	fmt.Println(a.move())
}

func main() {
	d := dog{}
	moveAnimal(d)
	f := fish{}
	moveAnimal(f)
	b := bird{}
	moveAnimal(b)
}
