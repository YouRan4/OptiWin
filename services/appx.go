//go:build windows

package services

import (
	"OptiWin/utils"
	"os/exec"
	"strings"
)

func ListAppxPackages() string {
	psScript := `
chcp 65001 >$null
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

Get-AppxPackage | ForEach-Object {
    $displayName = ""

    try {
        $entries = $_.GetAppListEntries()
        if ($entries -and $entries.Count -gt 0) {
            $displayName = $entries[0].DisplayInfo.DisplayName
        }
    } catch {}

    if (-not $displayName) {
        try {
            $manifest = Get-AppxPackageManifest -Package $_.PackageFullName -ErrorAction Stop
            if ($manifest.Package.Properties.DisplayName) {
                $displayName = $manifest.Package.Properties.DisplayName
            }
        } catch {}
    }

    if (-not $displayName -or $displayName.StartsWith("ms-resource:")) {
        try {
            $manifest = Get-AppxPackageManifest -Package $_.PackageFullName -ErrorAction Stop
            $resources = $manifest.Package.Resources
            if ($resources) {
                $lang = [System.Globalization.CultureInfo]::CurrentUICulture.Name
                $resource = $resources | Where-Object { $_.Language -eq $lang } | Select-Object -First 1
                if (-not $resource) {
                    $resource = $resources | Where-Object { $_.Language -eq "" } | Select-Object -First 1
                }
                if ($resource -and $resource.DisplayName -and -not $resource.DisplayName.StartsWith("ms-resource:")) {
                    $displayName = $resource.DisplayName
                }
            }
        } catch {}
    }

    if (-not $displayName -or $displayName.StartsWith("ms-resource:")) {
        try {
            $manifestPath = Join-Path $_.InstallLocation "AppxManifest.xml"
            if (Test-Path $manifestPath) {
                [xml]$xml = Get-Content $manifestPath -ErrorAction Stop
                if ($xml.Package.Properties.DisplayName) {
                    $displayName = $xml.Package.Properties.DisplayName
                }
            }
        } catch {}
    }

    if (-not $displayName -or $displayName.StartsWith("ms-resource:")) {
        $displayName = $_.Name
    }

    # 判断类型：用户应用或系统应用
    $pkgType = "user"
    if ($_.IsFramework -or $_.Name -match "^(Microsoft\.|Windows\.|WindowsCore\.)" -or $_.SignatureKind -eq "System") {
        $pkgType = "system"
    }

    [PSCustomObject]@{
        Name            = $_.Name
        PackageFullName = $_.PackageFullName
        Version         = [string]$_.Version
        DisplayName     = $displayName
        Type            = $pkgType
    }
} | ConvertTo-Json -Compress
`
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	utils.HideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return "[]"
	}
	result := strings.TrimSpace(string(out))
	if result == "" || result == "null" {
		return "[]"
	}
	return result
}

func UninstallAppx(fullName string) bool {
	psScript := "Remove-AppxPackage -Package '" + strings.ReplaceAll(fullName, "'", "''") + "'"
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", psScript)
	utils.HideWindow(cmd)
	return cmd.Run() == nil
}
