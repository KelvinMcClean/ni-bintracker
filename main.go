package main

import (
	"fmt"
	"log"

	"github.com/pelletier/go-toml"
)

func main() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	} else {

		// retrieve data directly

		homeId := config.Get("house.id").(int64)
		council := config.Get("council.name").(string)
		token := config.Get("calender.token").(string)
		fmt.Println("Home ID is", homeId, "and council is", council, "and token is", token)
	}
}
