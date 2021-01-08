package wslcli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type iPInfo struct {
	host  string
	local string
}

//ip route get to 1.0.0.0 | head -1 | awk '{print $7}'

// RunningDistros returns list of distros names running
func RunningDistros() ([]string, error) {
	cmd := exec.Command("wsl.exe", "-l", "-q", "--running")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	decoded, err := decodeOutput(out)
	if err != nil {
		return nil, errors.New("failed to decode output")
	}
	return strings.Split(decoded, "\r\n"), nil
}

// ListAll returns output for "wsl.exe -l -v"
func ListAll() (string, error) {
	cmd := exec.Command("wsl.exe", "-l", "-v")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("wsl -l -v failed: %w", err)
	}
	decoded, err := decodeOutput(out)
	if err != nil {
		return "", fmt.Errorf("failed to decode output: %w", err)
	}
	return decoded, nil
}

func GetIP(name string) (*iPInfo, error) {
	cmd := exec.Command("wsl.exe", "-d", name, "--", "ip", "route", "get", "to", "1.0.0.0")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("exec command(%s) failed: %w", strings.Join(cmd.Args, " "), err)
	}

	ipInfo := &iPInfo{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		// 10.0.0.1 via 172.25.176.1 dev eth0 src 172.25.185.133 uid 1000
		line := scanner.Text()
		line = strings.TrimSpace(line)
		info := strings.Split(line, " ")
		if len(info) < 9 {
			continue
		}
		hostIP := net.ParseIP(info[2])
		if hostIP == nil || hostIP.To4() == nil {
			continue
		}
		localIP := net.ParseIP(info[6])
		if localIP == nil || localIP.To4() == nil {
			continue
		}
		ipInfo.host = hostIP.String()
		ipInfo.local = localIP.String()
		return ipInfo, nil
	}

	return nil, fmt.Errorf("find ip failed")
}

func GetHostIP() (string, error) {
	// Get-NetIPAddress | Where-Object -FilterScript { $_.InterfaceAlias -eq "vEthernet (WSL)" } | Format-Table -HideTableHeaders -Wrap ifIndex,InterfaceAlias,IPAddress,AddressFamily,PrefixLength
	// Get-NetIPAddress -InterfaceAlias "vEthernet (WSL)" | Format-Table -HideTableHeaders -Wrap ifIndex,InterfaceAlias,IPAddress,AddressFamily,PrefixLength"
	cmd := exec.Command(
		"powershell.exe",
		"-Command",
		"Get-NetIPAddress -AddressFamily IPv4 -InterfaceAlias \"vEthernet (WSL)\" | Format-Table -HideTableHeaders -Wrap ifIndex,IPAddress,AddressFamily,PrefixLength",
	)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("wsl -l -v failed: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		var (
			fid          string
			ip           string
			nettype      string
			PrefixLength int
		)
		n, _ := fmt.Sscanf(line, "%s %s %s %d", &fid, &ip, &nettype, &PrefixLength)
		if n != 4 {
			continue
		}
		return ip, nil
	}
	return "", nil
}

// Get-NetIPAddress | Where-Object -FilterScript { $_.InterfaceAlias -eq "vEthernet (WSL)" }

func decodeOutput(raw []byte) (string, error) {
	win16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16le.NewDecoder())
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)
	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
