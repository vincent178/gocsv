package main

import (
	"fmt"
	"os"

	"github.com/vincent178/gocsv"
)

type Foo struct {
	Name   string
	Count  uint // speicify type other than string
	Enable bool `csv:"is_enable"` // speicify csv tag which used to mapping csv header
}

func main() {
	f, err := os.Open("./data.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ret, err := gocsv.Read[Foo](f)

	if err != nil {
		panic(err)
	}

	for _, r := range ret {
		fmt.Printf("%+v\n", r)
	}
}
