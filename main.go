package main

import (
	"fmt"

	"time-series-graphs/jozko"
)

func main() {
	mojAvto := jozko.Auto{
		Kolesa: 0,
		Okna:   0,
	}

	fmt.Println(mojAvto.Kolesa)
	jozko.F()
}
