package main

import (
	"fmt"

	"github.com/lukas-zx/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("lukas")
	fmt.Println(cfg)
	newCfg := config.Read()
	fmt.Println(newCfg)
}
