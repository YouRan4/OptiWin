//go:build windows

package utils

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

const th32csSnapProcess = 0x00000002

type processEntry32 struct {
	Size            uint32
	CntUsage        uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	CntThreads      uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [260]uint16
}

// TerminateProcesses 按映像名精确结束进程（不区分大小写），返回命中数量。
// 直接调 TerminateProcess，不经过 taskkill.exe，避免命令行特征被安全软件拦截；
// 例：TerminateProcesses("msedge.exe") 不会碰 msedgewebview2.exe。
func TerminateProcesses(imageName string) (int, error) {
	pids, err := processesByName(imageName)
	if err != nil {
		return 0, err
	}
	var errs []error
	for _, pid := range pids {
		if e := terminateProcess(pid); e != nil {
			errs = append(errs, fmt.Errorf("pid %d: %w", pid, e))
		}
	}
	return len(pids), errors.Join(errs...)
}

// CountProcesses 返回当前正在运行的指定映像名进程数（不区分大小写）
func CountProcesses(imageName string) (int, error) {
	pids, err := processesByName(imageName)
	return len(pids), err
}

func processesByName(imageName string) ([]uint32, error) {
	snap, _, err := kernel32.NewProc("CreateToolhelp32Snapshot").Call(th32csSnapProcess, 0)
	if snap == uintptr(syscall.InvalidHandle) {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot: %w", err)
	}
	defer kernel32.NewProc("CloseHandle").Call(snap)

	self, _, _ := kernel32.NewProc("GetCurrentProcessId").Call()

	var (
		entry processEntry32
		pids  []uint32
	)
	entry.Size = uint32(unsafe.Sizeof(entry))
	for ok, _, _ := kernel32.NewProc("Process32FirstW").Call(snap, uintptr(unsafe.Pointer(&entry))); ok != 0; ok, _, _ = kernel32.NewProc("Process32NextW").Call(snap, uintptr(unsafe.Pointer(&entry))) {
		if entry.ProcessID == 0 || uintptr(entry.ProcessID) == self {
			continue
		}
		if strings.EqualFold(syscall.UTF16ToString(entry.ExeFile[:]), imageName) {
			pids = append(pids, entry.ProcessID)
		}
	}
	return pids, nil
}
