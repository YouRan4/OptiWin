//go:build windows

package services

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strings"
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows/registry"
	"OptiWin/utils"
)

type systemInfo struct {
	OS    string `json:"os"`
	Build string `json:"build"`
	CPU   string `json:"cpu"`
	RAM   string `json:"ram"`
	IPv4  string `json:"ipv4"`
	IPv6  string `json:"ipv6"`
}

type memoryStatusEx struct {
	dwLength           uint32
	dwMemoryLoad       uint32
	ullTotalPhys       uint64
	ullAvailPhys       uint64
	ullTotalPageFile   uint64
	ullAvailPageFile   uint64
	ullTotalVirtual    uint64
	ullAvailVirtual    uint64
	ullExtendedVirtual uint64
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var procGlobalMemory = kernel32.NewProc("GlobalMemoryStatusEx")

func GetSystemInfo() string {
	info := systemInfo{
		OS:   getOS(),
		Build: getBuild(),
		CPU:  getCPU(),
		RAM:  getRAM(),
		IPv4: getIPv4(),
		IPv6: getIPv6(),
	}
	result, _ := json.Marshal(info)
	return string(result)
}

func getOS() string {
	name, _, _ := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "ProductName")
	ver, _, _ := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "DisplayVersion")
	build, _, _ := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "CurrentBuild")

	// Windows 11 开始 build >= 22000，但 ProductName 依然写 "Windows 10"
	if build >= "22000" {
		name = strings.Replace(name, "Windows 10", "Windows 11", 1)
	}

	if name != "" && ver != "" {
		return name + " " + ver
	}
	return name
}

func getBuild() string {
	b, _, _ := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "CurrentBuild")
	ubr, _ := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "UBR")
	if ubr > 0 {
		return fmt.Sprintf("%s.%d", b, ubr)
	}
	return b
}

func getCPU() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	name, _, _ := k.GetStringValue("ProcessorNameString")
	return strings.TrimSpace(name)
}

func getRAM() string {
	var state memoryStatusEx
	state.dwLength = uint32(unsafe.Sizeof(state))
	ret, _, _ := procGlobalMemory.Call(uintptr(unsafe.Pointer(&state)))
	if ret != 0 {
		gb := float64(state.ullTotalPhys) / (1024 * 1024 * 1024)
		return fmt.Sprintf("%.0f GB", math.Ceil(gb))
	}
	return ""
}

func getIPv4() string {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

func getIPv6() string {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ipv6 := ipnet.IP.To16(); ipv6 != nil && ipnet.IP.To4() == nil {
			return ipv6.String()
		}
	}
	return ""
}
