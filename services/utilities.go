//go:build windows
package services

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"
	"golang.org/x/sys/windows/registry"
	"OptiWin/utils"
)

func GetHibernateStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Power`, "HibernateEnabled")
	return err == nil && v == 1
}

func EnableHibernate() bool {
	var v byte = 1
	ret, _, _ := procCallNtPower.Call(systemReserveHiberFile,
		uintptr(unsafe.Pointer(&v)), 1, 0, 0)
	if ret != 0 {
		return false
	}
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Power`, "HibernateEnabled", 1)
	return true
}

func DisableHibernate() bool {
	var v byte = 0
	ret, _, _ := procCallNtPower.Call(systemReserveHiberFile,
		uintptr(unsafe.Pointer(&v)), 1, 0, 0)
	if ret != 0 {
		return false
	}
	utils.RegSetDWord(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Power`, "HibernateEnabled", 0)
	return true
}

func GetFastStartupStatus() bool {
	v, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Power`, "HiberbootEnabled")
	return err == nil && v == 1
}

func EnableFastStartup() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Power`, "HiberbootEnabled", 1)
	return true
}

func DisableFastStartup() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Power`, "HiberbootEnabled", 0)
	return true
}

func GetPhotoViewerStatus() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\.jpg\OpenWithProgids`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	vals, err := k.ReadValueNames(0)
	if err != nil {
		return false
	}
	for _, v := range vals {
		if strings.Contains(v, "PhotoViewer") {
			return true
		}
	}
	return false
}

func EnablePhotoViewer() bool {
	for _, ext := range imgExts {
		k, _, _ := registry.CreateKey(registry.CURRENT_USER,
			`Software\Classes\`+ext+`\OpenWithProgids`, registry.SET_VALUE)
		k.SetStringValue("PhotoViewer.FileAssoc.Tiff", "")
		k.Close()
	}
	return true
}

func DisablePhotoViewer() bool {
	for _, ext := range imgExts {
		k, err := registry.OpenKey(registry.CURRENT_USER,
			`Software\Classes\`+ext+`\OpenWithProgids`, registry.SET_VALUE)
		if err != nil {
			continue
		}
		k.DeleteValue("PhotoViewer.FileAssoc.Tiff")
		k.Close()
	}
	return true
}

func UninstallEdge() string {
	exec.Command("taskkill", "/f", "/im", "msedge.exe").Run()

	matches, err := filepath.Glob(`C:\Program Files (x86)\Microsoft\Edge\Application\*\Installer\setup.exe`)
	if err != nil || len(matches) == 0 {
		return "未找到 Edge 安装路径"
	}

	cmd := exec.Command(matches[0], "--uninstall", "--force-uninstall", "--system-level")
	utils.HideWindow(cmd)
	cmd.Run()

	edgeAppDir := filepath.Dir(filepath.Dir(matches[0]))
	os.RemoveAll(edgeAppDir)

	shortcuts := []string{
		os.Getenv("PUBLIC") + `\Desktop\Microsoft Edge.lnk`,
		os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Microsoft Edge.lnk`,
	}
	for _, s := range shortcuts {
		os.Remove(s)
	}

	return "Edge 已卸载（保留 EdgeUpdate/EdgeCore，WebView2 不受影响）"
}

func GetWebView2Version() string {
	guids := []string{
		`Software\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
		`SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`,
	}
	for _, g := range guids {
		v, _, err := utils.RegReadString(registry.LOCAL_MACHINE, g, "pv")
		if err == nil && v != "" {
			return v
		}
	}
	return "未安装"
}

func InstallWebView2() string {
	tmp := os.TempDir() + `\MicrosoftEdgeWebView2Setup.exe`
	url := "https://go.microsoft.com/fwlink/p/?LinkId=2124703"

	out, err := utils.GetHTTPClient().Get(url)
	if err != nil {
		return "下载失败: " + err.Error()
	}
	defer out.Body.Close()

	f, err := os.Create(tmp)
	if err != nil {
		return "创建文件失败: " + err.Error()
	}
	_, err = io.Copy(f, out.Body)
	f.Close()
	if err != nil {
		return "下载不完整: " + err.Error()
	}

	cmd := exec.Command(tmp, "/silent", "/install")
	utils.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return "安装失败: " + err.Error()
	}
	return "WebView2 安装/升级完成"
}

func SetSafeBoot(mode string) bool {
	utils.RunHide("bcdedit", "/deletevalue", "{current}", "safeboot")
	if mode == "minimal" {
		utils.RunHide("bcdedit", "/set", "{current}", "safeboot", "minimal")
	} else if mode == "network" {
		utils.RunHide("bcdedit", "/set", "{current}", "safeboot", "network")
	}
	return true
}

func RebootSystem() bool {
	cmd := exec.Command("shutdown", "/r", "/t", "0")
	utils.HideWindow(cmd)
	return cmd.Run() == nil
}

func RebootToBios() bool {
	cmd := exec.Command("shutdown", "/r", "/fw", "/t", "0")
	utils.HideWindow(cmd)
	return cmd.Run() == nil
}

// UpdateKGL 写入 KGL 数据到注册表
func UpdateKGL() string {
	data, errMsg := FetchKGL()
	if data == nil {
		return errMsg
	}

	kglKey, _, _ := registry.CreateKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\KGL\OneSettings`, registry.SET_VALUE)
	defer kglKey.Close()

	kglKey.SetStringValue("ActivateOnUpdate", fmt.Sprintf("%v", data.ActivateOnUpdate))
	kglKey.SetStringValue("Hash", data.Hash)
	kglKey.SetStringValue("URI", data.URI)
	kglKey.SetStringValue("Version", data.Version)
	kglKey.SetDWordValue("VersionCheckTimeout", uint32(data.VersionCheckTimeout))

	userKey, _, _ := registry.CreateKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\GameDVR`, registry.SET_VALUE)
	defer userKey.Close()

	userKey.SetStringValue("KGLRevision", data.Version)
	userKey.SetStringValue("KGLToGCSUpdatedRevision", data.Version)
	return "KGL 更新成功"
}
