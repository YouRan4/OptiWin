//go:build windows

package services

import (
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows"
)

var powrProf = syscall.NewLazyDLL("powrprof.dll")

var (
	procPowerGetActiveScheme      = powrProf.NewProc("PowerGetActiveScheme")
	procPowerSetActiveScheme      = powrProf.NewProc("PowerSetActiveScheme")
	procPowerDuplicatePowerScheme = powrProf.NewProc("PowerDuplicatePowerScheme")
	procPowerReadACValueIndex     = powrProf.NewProc("PowerReadACValueIndex")
	procPowerWriteACValueIndex    = powrProf.NewProc("PowerWriteACValueIndex")
)

var (
	ultimateGuid = windows.GUID{Data1: 0xe9a42b02, Data2: 0xd5df, Data3: 0x448d, Data4: [8]byte{0xaa, 0x00, 0x03, 0xf1, 0x47, 0x49, 0xeb, 0x61}}
	balancedGuid = windows.GUID{Data1: 0x381b4222, Data2: 0xf694, Data3: 0x41f0, Data4: [8]byte{0x96, 0x85, 0xff, 0x5b, 0xb2, 0x60, 0xdf, 0x2e}}
	subProcGuid  = windows.GUID{Data1: 0x54533251, Data2: 0x82be, Data3: 0x4824, Data4: [8]byte{0x96, 0xc1, 0x47, 0xb6, 0x0b, 0x74, 0x0d, 0x00}}
	cStateGuid   = windows.GUID{Data1: 0x0cc5b647, Data2: 0xc1df, Data3: 0x4637, Data4: [8]byte{0x89, 0x1a, 0xde, 0xc3, 0x5c, 0x31, 0x85, 0x83}}
)

func getActiveSchemeGuid() *windows.GUID {
	var p *windows.GUID
	ret, _, _ := procPowerGetActiveScheme.Call(0, uintptr(unsafe.Pointer(&p)))
	if ret != 0 {
		return nil
	}
	return p
}

func freeGuid(p *windows.GUID) {
	windows.LocalFree(windows.Handle(unsafe.Pointer(p)))
}

func cstateRead(guid *windows.GUID) (uint32, bool) {
	var val uint32
	ret, _, _ := procPowerReadACValueIndex.Call(0, uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(&subProcGuid)), uintptr(unsafe.Pointer(&cStateGuid)), uintptr(unsafe.Pointer(&val)))
	return val, ret == 0
}

func cstateWrite(guid *windows.GUID, val uint32) {
	procPowerWriteACValueIndex.Call(0, uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(&subProcGuid)), uintptr(unsafe.Pointer(&cStateGuid)), uintptr(val))
}

var procCallNtPower = powrProf.NewProc("CallNtPowerInformation")

const systemReserveHiberFile = 10

var imgExts = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif", ".tiff", ".tif", ".ico", ".webp"}
