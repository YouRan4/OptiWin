//go:build windows

package utils

import (
	_ "embed"
)

var (
	RemoveGameBarScript      []byte
	RestoreGameBarScript     []byte
	DisableDefenderScript    []byte
	RestoreDefenderScript    []byte
	EnableTaskManagerScript  []byte
	DisableTaskManagerScript []byte
)
