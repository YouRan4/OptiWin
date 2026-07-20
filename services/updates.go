//go:build windows
package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"golang.org/x/sys/windows/registry"
	"OptiWin/utils"
)

func GetUpdateChannel() string {
	branch, _, _ := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsSelfHost\Applicability`,
		"BranchName")
	ring, _, _ := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsSelfHost\Applicability`,
		"Ring")

	switch {
	case ring == "" || ring == "Retail":
		return "retail"
	case branch == "ReleasePreview":
		return "ReleasePreview"
	case branch == "Beta":
		return "Beta"
	case branch == "Dev":
		return "Dev"
	case branch == "Canary":
		return "Canary"
	default:
		return "retail"
	}
}

func SetUpdateChannel(channel string) bool {
	path := `SOFTWARE\Microsoft\WindowsSelfHost\Applicability`
	switch channel {
	case "retail":
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "Ring", "Retail")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "ContentType", "Mainline")
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.SET_VALUE)
		if err == nil {
			k.DeleteValue("BranchName")
			k.Close()
		}
	case "ReleasePreview":
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "BranchName", "ReleasePreview")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "Ring", "External")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "ContentType", "Skip")
	case "Beta":
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "BranchName", "Beta")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "Ring", "External")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "ContentType", "Skip")
	case "Dev":
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "BranchName", "Dev")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "Ring", "External")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "ContentType", "Skip")
	case "Canary":
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "BranchName", "Canary")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "Ring", "WIS")
		utils.RegWriteString(registry.LOCAL_MACHINE, path, "ContentType", "Skip")
	default:
		return false
	}
	return true
}

func GetPauseUpdatesStatus() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	_, _, err = k.GetStringValue("PauseUpdatesExpiryTime")
	return err == nil
}

func EnablePauseUpdates() bool {
	k, _, _ := registry.CreateKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.SET_VALUE)
	defer k.Close()

	k.SetDWordValue("FlightSettingsMaxPauseDays", 37200)
	k.SetStringValue("PauseFeatureUpdatesStartTime", "2023-08-17T12:47:51Z")
	k.SetStringValue("PauseFeatureUpdatesEndTime", "2126-01-01T00:00:00Z")
	k.SetStringValue("PauseQualityUpdatesStartTime", "2023-08-17T12:47:51Z")
	k.SetStringValue("PauseQualityUpdatesEndTime", "2126-01-01T00:00:00Z")
	k.SetStringValue("PauseUpdatesStartTime", "2023-08-17T12:47:51Z")
	k.SetStringValue("PauseUpdatesExpiryTime", "2126-01-01T00:00:00Z")
	return true
}

func DisablePauseUpdates() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.SET_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	for _, v := range []string{
		"FlightSettingsMaxPauseDays",
		"PauseFeatureUpdatesStartTime",
		"PauseFeatureUpdatesEndTime",
		"PauseQualityUpdatesStartTime",
		"PauseQualityUpdatesEndTime",
		"PauseUpdatesStartTime",
		"PauseUpdatesExpiryTime",
	} {
		k.DeleteValue(v)
	}
	return true
}

func GetVisibilityStatus() bool {
	val, _, err := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"SettingsPageVisibility")
	if err != nil || val == "" {
		return false
	}
	for i := 0; i+18 <= len(val); i++ {
		if val[i:i+18] == "hide:windowsupdate" {
			return true
		}
	}
	return false
}

func EnableVisibility() bool {
	existing, _, err := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"SettingsPageVisibility")
	if err != nil {
		existing = ""
	}
	if existing != "" {
		existing += ";hide:windowsupdate"
	} else {
		existing = "hide:windowsupdate"
	}
	cmd := exec.Command("reg", "add", `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"/v", "SettingsPageVisibility", "/t", "REG_SZ", "/d", existing, "/f")
	utils.HideWindow(cmd)
	if _, err := cmd.CombinedOutput(); err != nil {
		return false
	}
	return true
}

func DisableVisibility() bool {
	existing, _, err := utils.RegReadString(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"SettingsPageVisibility")
	if err != nil {
		return true
	}

	var parts []string
	start := 0
	for start < len(existing) {
		end := start
		for end < len(existing) && existing[end] != ';' {
			end++
		}
		part := existing[start:end]
		if part != "hide:windowsupdate" && part != "" {
			parts = append(parts, part)
		}
		start = end + 1
	}

	if len(parts) == 0 {
		cmd := exec.Command("reg", "delete", `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
			"/v", "SettingsPageVisibility", "/f")
		utils.HideWindow(cmd)
		if err := cmd.Run(); err != nil {
			return false
		}
		return true
	}

	newVal := parts[0]
	for i := 1; i < len(parts); i++ {
		newVal += ";" + parts[i]
	}
	cmd := exec.Command("reg", "add", `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		"/v", "SettingsPageVisibility", "/t", "REG_SZ", "/d", newVal, "/f")
	utils.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func GetDriverUpdatesStatus() bool {
	val, err := utils.RegReadDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Device Metadata`,
		"PreventDeviceMetadataFromNetwork")
	if err != nil || val == 0 {
		return true
	}
	return false
}

func EnableDriverUpdates() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Device Metadata`,
		"PreventDeviceMetadataFromNetwork", 0)
	utils.RegDeleteValue(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate`,
		"ExcludeWUDriversInQualityUpdate")
	utils.RegDeleteValue(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`,
		"ExcludeWUDriversInQualityUpdate")
	return true
}

func DisableDriverUpdates() bool {
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Device Metadata`,
		"PreventDeviceMetadataFromNetwork", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate`,
		"ExcludeWUDriversInQualityUpdate", 1)
	utils.RegSetDWord(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`,
		"ExcludeWUDriversInQualityUpdate", 1)
	return true
}

// --- 保留原有内容 ---

func UpdateCertificates() string {
	b := make([]byte, 8)
	rand.Read(b)
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("certs_%x.sst", b))
	defer os.Remove(tmpFile)

	cmd := exec.Command("CertUtil", "-generateSSTFromWU", tmpFile)
	utils.HideWindow(cmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Sprintf("下载证书失败: %v\n输出: %s", err, string(out))
	}

	stat, err := os.Stat(tmpFile)
	if err != nil || stat.Size() == 0 {
		return "证书文件无效或为空"
	}

	importRoot := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Import-Certificate -FilePath "%s" -CertStoreLocation Cert:\LocalMachine\Root`, tmpFile))
	utils.HideWindow(importRoot)
	if out, err := importRoot.CombinedOutput(); err != nil {
		return fmt.Sprintf("导入根证书失败: %v\n输出: %s", err, string(out))
	}

	importAuth := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Import-Certificate -FilePath "%s" -CertStoreLocation Cert:\LocalMachine\AuthRoot`, tmpFile))
	utils.HideWindow(importAuth)
	if out, err := importAuth.CombinedOutput(); err != nil {
		return fmt.Sprintf("导入颁发机构证书失败: %v\n输出: %s", err, string(out))
	}

	return "证书更新成功"
}

type KGLData struct {
	Version             string
	ActivateOnUpdate    bool
	Hash                string
	URI                 string
	VersionCheckTimeout int
}

func FetchKGL() (*KGLData, string) {
	client := utils.GetHTTPClient()
	resp, err := client.Get("https://settings.data.microsoft.com/settings/v3.0/xbox/knowngamelist")
	if err != nil {
		return nil, fmt.Sprintf("KGL 请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Sprintf("KGL 读取响应失败: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Sprintf("KGL JSON 解析失败: %v", err)
	}

	settingsRaw, ok := result["settings"]
	if !ok {
		return nil, "KGL 响应中缺少 settings 字段"
	}

	var raw struct {
		Version             string `json:"version"`
		ActivateOnUpdate    bool   `json:"activateOnUpdate"`
		Hash                string `json:"hash"`
		URI                 string `json:"uri"`
		VersionCheckTimeout int    `json:"versionCheckTimeout"`
	}
	if err := json.Unmarshal(settingsRaw, &raw); err != nil {
		return nil, fmt.Sprintf("KGL settings 解析失败: %v", err)
	}

	return &KGLData{
		Version:             raw.Version,
		ActivateOnUpdate:    raw.ActivateOnUpdate,
		Hash:                raw.Hash,
		URI:                 raw.URI,
		VersionCheckTimeout: raw.VersionCheckTimeout,
	}, ""
}
