package main

import (
	"fmt"
	"zxin.org/demo/pkg/wslcli"
)

func main() {

	//h, err := hostsapi.CreateAPI()
	//
	//fmt.Println(err)
	//fmt.Println(h)
	//h.Parse()

	fmt.Println(wslcli.GetHostIP())
}
