package main

import (
	"strconv"
	"strings"
)

const CurrentVersion = "v1.3.1"

func parseVersion(v string) (major, minor, patch int) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) >= 1 {
		major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		patchStr := parts[2]
		if idx := strings.Index(patchStr, "-"); idx >= 0 {
			patchStr = patchStr[:idx]
		}
		patch, _ = strconv.Atoi(patchStr)
	}
	return
}

func compareVersion(current, remote string) int {
	cm, cn, cp := parseVersion(current)
	rm, rn, rp := parseVersion(remote)
	if cm != rm {
		if cm < rm {
			return -1
		}
		return 1
	}
	if cn != rn {
		if cn < rn {
			return -1
		}
		return 1
	}
	if cp != rp {
		if cp < rp {
			return -1
		}
		return 1
	}
	return 0
}
