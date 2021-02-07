package main

import (
	"fmt"
	"zxin.org/demo/pkg/hostsapi"
)

type Config struct {
}

func main() {
	h, _ := hostsapi.CreateAPI()
	h.Parse()
	//fmt.Println(h.GetHostLines())

	fmt.Println("===========================")

	lines := h.GetHostLines()
	for ix := 0; ix < len(lines); ix++ {
		fmt.Println(lines[ix])
	}

	fmt.Println("===========================")

	blocks := h.GetBlocks()
	for ix := 0; ix < len(blocks); ix++ {
		fmt.Println(blocks[ix])
	}

	hosts := make(map[int]*hostsapi.HostNameLine)
	hosts[0] = &hostsapi.HostNameLine{
		Address:  "172.18.156.43",
		Hostname: "ubuntu2004.wsl",
	}

	h.UpdateBlock("Ubuntu-20.04", hosts)

	//fmt.Println(wslcli.ListAll())
	//fmt.Println(wslcli.RunningDistros())
	//fmt.Println(wslcli.GetIP("Ubuntu-20.04"))
	//fmt.Println(wslcli.GetHostIP())
}
