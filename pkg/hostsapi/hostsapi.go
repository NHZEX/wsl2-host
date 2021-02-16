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
	LineFlagBlock    = 1
	LineFlagBegin    = 2
	LineFlagEnd      = 4
	LineFlagHostname = 8
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
	lines map[int]*HostnameLine
}

type HostnameLine struct {
	Address  string
	Hostname string
	idx      string
}

type HostsAPI struct {
	file        *os.File
	lines       map[int]*HostLine
	finder      int
	tmpBlock    *hostBlockMeta
	blockPos    map[int]string
	blocks      map[string]*hostBlockMeta
	lineCapture *regexp.Regexp
	lineCut     *regexp.Regexp
}

func CreateAPI() (*HostsAPI, error) {
	// 初始化正则表达式
	lineCapture := regexp.MustCompile("####@ (Begin|End) Hosts \\((.+?)\\) @####")
	lineCut := regexp.MustCompile("^(.+?)\\s+(.+)$")

	systemRoot := os.Getenv("SystemRoot")
	f, err := os.Open(path.Join(systemRoot, hostsPath))
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	h := &HostsAPI{
		file:        f,
		lines:       make(map[int]*HostLine),
		blockPos:    make(map[int]string),
		blocks:      make(map[string]*hostBlockMeta),
		lineCapture: lineCapture,
		lineCut:     lineCut,
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
			h.addBlockLine(line)
		}
		return
	}
	regexpResult := h.lineCapture.FindStringSubmatch(content)
	if len(regexpResult) == 0 && h.tmpBlock != nil {
		h.addBlockLine(line)
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
				lines: make(map[int]*HostnameLine),
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
			h.blocks[name] = h.tmpBlock
			h.blockPos[h.tmpBlock.begin] = name
			h.tmpBlock = nil
		}
	}
}

func (h *HostsAPI) addBlockLine(line *HostLine) {
	line.Flags = LineFlagBlock
	regexpResult := h.lineCut.FindStringSubmatch(line.Content)
	if len(regexpResult) != 3 {
		return
	}
	h.tmpBlock.lines[len(h.tmpBlock.lines)] = &HostnameLine{
		Address:  regexpResult[1],
		Hostname: regexpResult[2],
	}
	line.Flags = LineFlagBlock | LineFlagHostname
}

func (h *HostsAPI) GetHostLines() map[int]*HostLine {
	return h.lines
}

func (h *HostsAPI) GetBlocks() map[string]*hostBlockMeta {
	return h.blocks
}

func (h *HostsAPI) ChangeBlockHostnames(wslName string, hosts map[int]*HostnameLine) {
	if _, has := h.blocks[wslName]; !has {
		h.blocks[wslName] = &hostBlockMeta{
			begin: 0,
			end:   0,
			name:  wslName,
			lines: hosts,
		}
	} else {
		h.blocks[wslName].lines = hosts
	}
}

func (h *HostsAPI) AddBlockHostname(name string, line HostnameLine) {
	// 查找块
	// 添加行
}

func (h *HostsAPI) HasBlock(wslName string) bool {
	if _, has := h.blocks[wslName]; !has {
		return false
	}
	return true
}

func (h *HostsAPI) HasBlockHostname(wslName string, hostname string) bool {
	block, has := h.blocks[wslName]
	if !has {
		return false
	}

	for _, line := range block.lines {
		if line.Hostname == hostname {
			return true
		}
	}

	return false
}

func (h *HostsAPI) build() string {
	content := ""

	for idx := 0; idx < len(h.lines); idx++ {
		line := h.lines[idx]

		if line.Flags&LineFlagBegin > 0 {
			blockName, has := h.blockPos[line.idx]
			if !has {
				continue
			}

			block := h.blocks[blockName]
			content += fmt.Sprintf("####@ Begin Hosts (%s) @####\r\n", block.name)

			for idx := 0; idx < len(block.lines); idx++ {
				line := block.lines[idx]
				content += fmt.Sprintf("%s %s\r\n", line.Address, line.Hostname)
			}
			content += fmt.Sprintf("####@ End Hosts (%s) @####\r\n", block.name)
		} else if line.Flags&LineFlagBlock > 0 {
			continue
		} else {
			//fmt.Println(line)
			content += line.Content
			content += "\r\n"
		}
	}

	// 处理新增快
	for _, block := range h.blocks {
		if _, has := h.blockPos[block.begin]; has {
			continue
		}

		content += fmt.Sprintf("\r\n####@ Begin Hosts (%s) @####\r\n", block.name)

		for idx := 0; idx < len(block.lines); idx++ {
			line := block.lines[idx]
			content += fmt.Sprintf("%s %s\r\n", line.Address, line.Hostname)
		}
		content += fmt.Sprintf("####@ End Hosts (%s) @####\r\n", block.name)
	}

	return content
}

func (h *HostsAPI) Save() bool {
	content := h.build()

	// todo 调试
	fmt.Println(content)

	return true
}
