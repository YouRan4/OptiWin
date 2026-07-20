package main

import (
	"OptiWin/services"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// 重启 Windows 资源管理器
func (a *App) RestartExplorer() string {
	if runtime.GOOS != "windows" {
		return "仅支持 Windows"
	}
	cmd := exec.Command("taskkill", "/f", "/im", "explorer.exe")
	hideWindow(cmd)
	cmd.Run()

	go func() {
		cmd2 := exec.Command("explorer.exe")
		hideWindow(cmd2)
		cmd2.Start()
	}()
	return "资源管理器已重启"
}

// 性能
func (a *App) TogglePowerPlan(on bool) bool   { return services.SetToggle("power_plan", on) }
func (a *App) ToggleSuperfetch(on bool) bool  { return services.SetToggle("superfetch", on) }
func (a *App) ToggleMemCompress(on bool) bool { return services.SetToggle("mem_compress", on) }

// 个性化
func (a *App) ToggleNotification(mode string) bool { return services.SetNotification(mode) }

// 实用工具
func (a *App) ToggleHibernate(on bool) bool { return services.SetToggle("hibernate", on) }

// 更新
func (a *App) UpdateCertificates() string       { return services.UpdateCertificates() }
func (a *App) ToggleDriverUpdates(on bool) bool { return services.SetToggle("driver_updates", on) }

func (a *App) GetCurrentVersion() string { return CurrentVersion }

// CheckUpdate 检查 GitHub Release 最新版本
// 返回 JSON: {"version":"v1.0.1","body":"更新内容...","url":"..."} 或 "" (已是最新)
func (a *App) CheckUpdate() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/YouRan4/OptiWin/releases/latest", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.Unmarshal(b, &release); err != nil {
		return ""
	}

	if release.TagName == "" || release.TagName == CurrentVersion {
		return ""
	}

	cur := strings.TrimPrefix(CurrentVersion, "v")
	lat := strings.TrimPrefix(release.TagName, "v")
	if compareVer(cur, lat) >= 0 {
		return ""
	}

	result, _ := json.Marshal(map[string]string{
		"version": release.TagName,
		"body":    release.Body,
		"url":     release.HTMLURL,
	})
	return string(result)
}

func compareVer(a, b string) int {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		var ai, bi int
		if i < len(as) {
			ai = atoi(as[i])
		}
		if i < len(bs) {
			bi = atoi(bs[i])
		}
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}
	return 0
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
