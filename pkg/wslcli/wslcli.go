package wslcli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

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
