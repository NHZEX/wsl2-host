package hostsapi

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

const hostsPath = "System32\\drivers\\etc\\hosts"

type HostsAPI struct {
	file *os.File
}

func CreateAPI() (*HostsAPI, error) {
	systemRoot := os.Getenv("SystemRoot")
	f, err := os.Open(path.Join(systemRoot, hostsPath))
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	h := &HostsAPI{
		file: f,
	}
	return h, nil
}

func (h *HostsAPI) Parse() {
	scanner := bufio.NewScanner(h.file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}
