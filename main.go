package main

import (
	"fmt"

	"github.com/kevvarlar/gator/internal/config"
)
func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	err = cfg.SetUser("kevvarlar")
	if err != nil {
		fmt.Println(err)
	}
	result, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}