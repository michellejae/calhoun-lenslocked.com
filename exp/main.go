package main

import (
	"fmt"

	"gitlab.com/michellejae/lenslocked.com/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
