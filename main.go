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

	//blocks := h.GetBlocks()
	//for name, block := range blocks {
	//	fmt.Println(name, block)
	//}

	//fmt.Println(h.HasBlock("Ubuntu-20.04"))
	//fmt.Println(h.HasBlockHostname("Ubuntu-20.04", "ubuntu2004.wsl"))

	hosts := make(map[int]*hostsapi.HostnameLine)
	hosts[0] = &hostsapi.HostnameLine{
		Address:  "172.18.156.44",
		Hostname: "ubuntu2004.wsl",
	}
	hosts[1] = &hostsapi.HostnameLine{
		Address:  "172.18.156.43",
		Hostname: "wsl.local",
	}
	hosts[2] = &hostsapi.HostnameLine{
		Address:  "172.18.156.43",
		Hostname: "www.nxadmin.test",
	}
	hosts[3] = &hostsapi.HostnameLine{
		Address:  "172.18.156.43",
		Hostname: "www.sshotel.test",
	}
	h.ChangeBlockHostnames("Ubuntu-20.0411", hosts)
	h.Save()

	//fmt.Println(wslcli.ListAll())
	//fmt.Println(wslcli.RunningDistros())
	//fmt.Println(wslcli.GetIP("Ubuntu-20.04"))
	//fmt.Println(wslcli.GetHostIP())
}
