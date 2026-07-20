package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// UpdateCertificates 从 Windows Update 拉取根证书并导入系统
func UpdateCertificates() string {
	if runtime.GOOS != "windows" {
		return "仅支持 Windows"
	}

	// 生成临时文件名
	b := make([]byte, 8)
	rand.Read(b)
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("certs_%x.sst", b))
	defer os.Remove(tmpFile)

	// 调用 CertUtil 下载证书
	cmd := exec.Command("CertUtil", "-generateSSTFromWU", tmpFile)
	hideWindow(cmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Sprintf("下载证书失败: %v\n输出: %s", err, string(out))
	}

	// 检查文件是否有效
	stat, err := os.Stat(tmpFile)
	if err != nil || stat.Size() == 0 {
		return "证书文件无效或为空"
	}

	// 导入到受信任根证书存储
	importRoot := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Import-Certificate -FilePath "%s" -CertStoreLocation Cert:\LocalMachine\Root`, tmpFile))
	hideWindow(importRoot)
	if out, err := importRoot.CombinedOutput(); err != nil {
		return fmt.Sprintf("导入根证书失败: %v\n输出: %s", err, string(out))
	}

	// 导入到受信任颁发机构存储
	importAuth := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Import-Certificate -FilePath "%s" -CertStoreLocation Cert:\LocalMachine\AuthRoot`, tmpFile))
	hideWindow(importAuth)
	if out, err := importAuth.CombinedOutput(); err != nil {
		return fmt.Sprintf("导入颁发机构证书失败: %v\n输出: %s", err, string(out))
	}

	return "证书更新成功"
}

// KGLData 知识图谱库数据结构
type KGLData struct {
	Version              string
	ActivateOnUpdate     bool
	Hash                 string
	URI                  string
	VersionCheckTimeout  int
}

// FetchKGL 从微软 API 获取最新 KGL 数据
func FetchKGL() (*KGLData, string) {
	client := &http.Client{Timeout: 15 * time.Second}
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
		Version              string `json:"version"`
		ActivateOnUpdate     bool   `json:"activateOnUpdate"`
		Hash                 string `json:"hash"`
		URI                  string `json:"uri"`
		VersionCheckTimeout  int    `json:"versionCheckTimeout"`
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
