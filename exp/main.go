package main

import "fmt"

type Dog struct{}
type Cat struct{}

func (d Dog) Speak() {
	fmt.Println("wooof")
}

func (c Cat) Speak() {
	fmt.Println("meow")
}

type Husky struct {
	Speaker
}

type SpeakerPrefixer struct {
	Speaker
}

func (sp SpeakerPrefixer) Speak() {
	fmt.Print("Prefix:")
	sp.Speaker.Speak()
}

type Speaker interface {
	Speak()
}

func main() {
	h := Husky{Dog{}}
	h.Speak() // h.Dog.Speak() --> Dog

	// Husky doesn't care if we pass in Cat or Dog as long as the Speak method is on the Struct
	// cause we pass the interface Spaker into Huskey which just has a Speak method
	p := Husky{Cat{}}
	p.Speak() // p.Cat.Speak() --> Meow

	// something about chaining interfaces
	t := Husky{SpeakerPrefixer{Cat{}}}
	t.Speak() // --> Prefix: meow
}
