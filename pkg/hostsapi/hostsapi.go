package hostsapi

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

const hostsPath = "System32\\drivers\\etc\\hosts"

const (
	LineFlagBlock = 1
	LineFlagBegin = 2
	LineFlagEnd   = 4
	LineFlagHost  = 8
)

type HostLine struct {
	idx     int
	Content string
	Flags   int
}

type hostBlockMeta struct {
	begin int
	end   int
	name  string
	lines map[int]*HostNameLine
}

type HostNameLine struct {
	Address  string
	Hostname string
	idx      string
}

type HostsAPI struct {
	file        *os.File
	lines       map[int]*HostLine
	finder      int
	tmpBlock    *hostBlockMeta
	blocks      map[int]*hostBlockMeta
	lineCapture *regexp.Regexp
}

func CreateAPI() (*HostsAPI, error) {
	// 初始化正则表达式
	lineCapture := regexp.MustCompile("####@ (Begin|End) Hosts \\((.+?)\\) @####")

	systemRoot := os.Getenv("SystemRoot")
	f, err := os.Open(path.Join(systemRoot, hostsPath))
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	h := &HostsAPI{
		file:        f,
		lines:       make(map[int]*HostLine),
		blocks:      make(map[int]*hostBlockMeta),
		lineCapture: lineCapture,
	}
	return h, nil
}

func (h *HostsAPI) Parse() {
	scanner := bufio.NewScanner(h.file)
	idx := 0
	for scanner.Scan() {
		content := scanner.Text()
		line := &HostLine{
			idx:     idx,
			Content: content,
		}

		h.parseLine(line)

		h.lines[idx] = line
		idx++
	}
}

func (h *HostsAPI) parseLine(line *HostLine) {
	content := strings.TrimSpace(line.Content)
	if !strings.HasPrefix(content, "####@") || !strings.HasSuffix(content, "@####") {
		if h.tmpBlock == nil {
			line.Flags = 0
		} else {
			// todo 验证是否合法 host
			line.Flags = LineFlagBlock
		}
		return
	}
	regexpResult := h.lineCapture.FindStringSubmatch(content)
	if len(regexpResult) == 0 && h.tmpBlock != nil {
		line.Flags = LineFlagBlock
		return
	} else if len(regexpResult) != 3 {
		line.Flags = 0
		return
	}
	head := regexpResult[1]
	name := regexpResult[2]
	if h.tmpBlock == nil {
		if head == "Begin" {
			h.tmpBlock = &hostBlockMeta{
				begin: line.idx,
				end:   0,
				name:  name,
				lines: make(map[int]*HostNameLine),
			}
			line.Flags = LineFlagBlock | LineFlagBegin
		} else {
			line.Flags = 0
			return
		}
	} else {
		if head == "Begin" {
			line.Flags = 0
			h.tmpBlock = nil
			return
		} else if head == "End" {
			if name != h.tmpBlock.name {
				line.Flags = 0
				h.tmpBlock = nil
				return
			}
			line.Flags = LineFlagBlock | LineFlagEnd
			h.tmpBlock.end = line.idx
			h.blocks[len(h.blocks)] = h.tmpBlock
			h.tmpBlock = nil
		}
	}
}

func (h *HostsAPI) GetHostLines() map[int]*HostLine {
	return h.lines
}

func (h *HostsAPI) GetBlocks() map[int]*hostBlockMeta {
	return h.blocks
}

func (h *HostsAPI) UpdateBlock(name string, hosts map[int]*HostNameLine) {
	// 更新块
}

func (h *HostsAPI) AddBlockLine(name string, line HostNameLine) {
	// 查找块
	// 添加行
	// line := HostNameLine{
	// 	Address: "172.18.156.43",
	// 	Hostname: "ubuntu2004.wsl",
	// }
}
