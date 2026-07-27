//go:build windows

package tilauncher

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	tokenPrimary         = 1
	tokenAssignPrimary   = 0x0001
	tokenDuplicateAccess = 0x0002
	tokenAdjustPrivs     = 0x0020
	tokenQuery           = 0x0008
	maximumAllowed       = 0x02000000
	seDebugPrivilegeName = "SeDebugPrivilege"
	seImpersonateName    = "SeImpersonatePrivilege"
)

var (
	modadvapi32                  = windows.NewLazySystemDLL("advapi32.dll")
	procCreateProcessWithTokenW  = modadvapi32.NewProc("CreateProcessWithTokenW")
	procImpersonateLoggedOnUser  = modadvapi32.NewProc("ImpersonateLoggedOnUser")
	tiMutex                      sync.Mutex
)

// RunAsTrustedInstaller 以 TrustedInstaller 权限运行传入的命令
func RunAsTrustedInstaller(cmd *exec.Cmd) error {
	tiMutex.Lock()
	defer tiMutex.Unlock()

	if err := enablePrivileges(); err != nil {
		return err
	}

	tiPID, err := startTrustedInstallerService()
	if err != nil {
		return err
	}

	tiToken, err := getTokenViaImpersonation(tiPID)
	if err != nil {
		return err
	}
	defer tiToken.Close()

	var dupToken windows.Token
	if err := windows.DuplicateTokenEx(
		tiToken, maximumAllowed, &windows.SecurityAttributes{},
		windows.SecurityImpersonation, tokenPrimary, &dupToken,
	); err != nil {
		return fmt.Errorf("DuplicateTokenEx: %w", err)
	}
	defer dupToken.Close()

	cmdLine := buildCommandLine(cmd)
	if err := createProcessWithToken(dupToken, cmdLine, cmd.Dir); err != nil {
		return err
	}
	return nil
}

// --- 权限管理 ---

// enablePrivileges 启用 SeDebugPrivilege 和 SeImpersonatePrivilege
func enablePrivileges() error {
	var hToken windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), tokenAdjustPrivs|tokenQuery, &hToken); err != nil {
		return err
	}
	defer hToken.Close()

	var luidDebug, luidImp windows.LUID
	for _, pair := range [2]struct {
		name string
		out  *windows.LUID
	}{
		{seDebugPrivilegeName, &luidDebug},
		{seImpersonateName, &luidImp},
	} {
		p, _ := windows.UTF16PtrFromString(pair.name)
		if err := windows.LookupPrivilegeValue(nil, p, pair.out); err != nil {
			return fmt.Errorf("LookupPrivilegeValue(%s): %w", pair.name, err)
		}
	}

	type tpBuf struct {
		PrivilegeCount uint32
		Privileges     [2]windows.LUIDAndAttributes
	}
	tkp := tpBuf{PrivilegeCount: 2, Privileges: [2]windows.LUIDAndAttributes{
		{Luid: luidDebug, Attributes: 0x00000002},
		{Luid: luidImp, Attributes: 0x00000002},
	}}

	if err := windows.AdjustTokenPrivileges(hToken, false,
		(*windows.Tokenprivileges)(unsafe.Pointer(&tkp)), uint32(unsafe.Sizeof(tkp)), nil, nil); err != nil {
		return fmt.Errorf("AdjustTokenPrivileges: %w", err)
	}

	var returnLen uint32
	_ = windows.GetTokenInformation(hToken, windows.TokenPrivileges, nil, 0, &returnLen)
	buf := make([]byte, returnLen)
	if err := windows.GetTokenInformation(hToken, windows.TokenPrivileges, &buf[0], returnLen, &returnLen); err != nil {
		return err
	}
	tp := (*windows.Tokenprivileges)(unsafe.Pointer(&buf[0]))

	var hasDebug, hasImp bool
	for _, p := range tp.AllPrivileges() {
		if p.Luid == luidDebug && p.Attributes&0x00000002 != 0 {
			hasDebug = true
		}
		if p.Luid == luidImp && p.Attributes&0x00000002 != 0 {
			hasImp = true
		}
	}
	if !hasDebug {
		return fmt.Errorf("SeDebugPrivilege 未启用")
	}
	if !hasImp {
		return fmt.Errorf("SeImpersonatePrivilege 未启用")
	}
	return nil
}

// --- 服务管理 ---

// startTrustedInstallerService 启动 TrustedInstaller 服务并返回 PID
func startTrustedInstallerService() (uint32, error) {
	hSCM, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return 0, err
	}
	defer windows.CloseServiceHandle(hSCM)

	name, _ := windows.UTF16PtrFromString("TrustedInstaller")
	hSvc, err := windows.OpenService(hSCM, name, windows.SERVICE_START|windows.SERVICE_QUERY_STATUS)
	if err != nil {
		return 0, err
	}
	defer windows.CloseServiceHandle(hSvc)

	_ = windows.StartService(hSvc, 0, nil)

	for i := 0; i < 10; i++ {
		if pid := findSystemProcess("TrustedInstaller.exe", "servicing"); pid != 0 {
			return pid, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return 0, fmt.Errorf("未找到 TrustedInstaller 进程")
}

// findSystemProcess 查找系统目录下的指定进程，验证完整路径防止进程名欺骗
func findSystemProcess(name string, subdir string) uint32 {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(snapshot, &entry); err != nil {
		return 0
	}

	expectedPath := strings.ToLower(os.Getenv("SystemRoot") + "\\" + subdir + "\\" + name)
	for {
		if windows.UTF16ToString(entry.ExeFile[:]) == name {
			if pid := verifyProcessPath(entry.ProcessID, expectedPath); pid != 0 {
				return pid
			}
		}
		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}
	return 0
}

// verifyProcessPath 用 QueryFullProcessImageName 验证进程的完整路径
func verifyProcessPath(pid uint32, expectedPath string) uint32 {
	hProc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(hProc)

	var exePath [windows.MAX_PATH]uint16
	size := uint32(len(exePath))
	if err := windows.QueryFullProcessImageName(hProc, 0, &exePath[0], &size); err != nil {
		return 0
	}
	if strings.ToLower(windows.UTF16ToString(exePath[:size])) == expectedPath {
		return pid
	}
	return 0
}

// --- 令牌获取 ---

// getTokenViaImpersonation 通过模拟 SYSTEM 身份获取 TrustedInstaller 令牌
func getTokenViaImpersonation(tiPID uint32) (windows.Token, error) {
	winPID := findSystemProcess("winlogon.exe", "System32")
	if winPID == 0 {
		return 0, fmt.Errorf("未找到 winlogon.exe")
	}

	hWin, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, winPID)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(hWin)

	var hSysToken windows.Token
	if err := windows.OpenProcessToken(hWin, tokenQuery|tokenDuplicateAccess, &hSysToken); err != nil {
		return 0, err
	}
	defer hSysToken.Close()

	var sysPrimary windows.Token
	if err := windows.DuplicateTokenEx(hSysToken, windows.MAXIMUM_ALLOWED, nil,
		windows.SecurityImpersonation, tokenPrimary, &sysPrimary); err != nil {
		return 0, err
	}
	defer sysPrimary.Close()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if _, _, e := syscall.SyscallN(procImpersonateLoggedOnUser.Addr(), uintptr(sysPrimary)); e != 0 {
		return 0, fmt.Errorf("ImpersonateLoggedOnUser: %w", e)
	}
	defer windows.RevertToSelf()

	hTI, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, tiPID)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(hTI)

	var hTIToken windows.Token
	if err := windows.OpenProcessToken(hTI, tokenQuery|tokenDuplicateAccess|tokenAssignPrimary, &hTIToken); err != nil {
		return 0, err
	}
	return hTIToken, nil
}

// --- 命令行构建 ---

func buildCommandLine(cmd *exec.Cmd) string {
	var sb strings.Builder
	sb.WriteString(syscall.EscapeArg(cmd.Path))
	for _, arg := range cmd.Args[1:] {
		sb.WriteByte(' ')
		sb.WriteString(syscall.EscapeArg(arg))
	}
	return sb.String()
}

// --- 进程创建 ---

func createProcessWithToken(token windows.Token, cmdLine string, workDir string) error {
	cmdLinePtr, err := windows.UTF16PtrFromString(cmdLine)
	if err != nil {
		return err
	}

	var workDirPtr *uint16
	if workDir != "" {
		workDirPtr, _ = windows.UTF16PtrFromString(workDir)
	}

	var si windows.StartupInfo
	si.Cb = uint32(unsafe.Sizeof(si))

	var pi windows.ProcessInformation

	r1, _, e1 := syscall.SyscallN(
		procCreateProcessWithTokenW.Addr(),
		uintptr(token), 0, 0,
		uintptr(unsafe.Pointer(cmdLinePtr)),
		0x08000000, // CREATE_NO_WINDOW
		0, uintptr(unsafe.Pointer(workDirPtr)),
		uintptr(unsafe.Pointer(&si)), uintptr(unsafe.Pointer(&pi)),
	)
	if r1 == 0 {
		return fmt.Errorf("CreateProcessWithTokenW: %w", errnoErr(e1))
	}

	windows.WaitForSingleObject(pi.Process, windows.INFINITE)

	var exitCode uint32
	exitErr := windows.GetExitCodeProcess(pi.Process, &exitCode)
	windows.CloseHandle(pi.Process)
	windows.CloseHandle(pi.Thread)

	if exitErr != nil {
		return fmt.Errorf("GetExitCodeProcess: %w", exitErr)
	}
	if exitCode != 0 {
		return fmt.Errorf("进程异常退出，退出码: %d", exitCode)
	}
	return nil
}

func errnoErr(e syscall.Errno) error {
	if e != 0 {
		return e
	}
	return nil
}
