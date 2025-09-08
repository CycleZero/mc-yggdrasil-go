package main

import (
	"fmt"
	"github.com/CycleZero/mc-yggdrasil-go/utils"
)

func main() {
	name := "Poyuan233"
	pid, err := utils.NameUUIDFromBytes([]byte(name))
	if err != nil {
		panic(err)
	}
	fmt.Println(pid)

}
