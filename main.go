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

	fmt.Println(wslcli.ListAll())
	fmt.Println(wslcli.RunningDistros())
	fmt.Println(wslcli.GetIP("Ubuntu-20.04"))
	fmt.Println(wslcli.GetHostIP())
}
