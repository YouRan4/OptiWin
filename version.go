package main

import (
	"fmt"
	"strconv"
	"strings"
)

const CurrentVersion = "v1.8.1"

func parseVersion(v string) (major, minor, patch int, err error) {
	orig := v
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) >= 1 {
		major, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("无效版本号 %q: %w", orig, err)
		}
	}
	if len(parts) >= 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("无效版本号 %q: %w", orig, err)
		}
	}
	if len(parts) >= 3 {
		patchStr := parts[2]
		if idx := strings.Index(patchStr, "-"); idx >= 0 {
			patchStr = patchStr[:idx]
		}
		patch, err = strconv.Atoi(patchStr)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("无效版本号 %q: %w", orig, err)
		}
	}
	return major, minor, patch, nil
}

func compareVersion(current, remote string) (int, error) {
	cm, cn, cp, err := parseVersion(current)
	if err != nil {
		return 0, err
	}
	rm, rn, rp, err := parseVersion(remote)
	if err != nil {
		return 0, err
	}
	if cm != rm {
		if cm < rm {
			return -1, nil
		}
		return 1, nil
	}
	if cn != rn {
		if cn < rn {
			return -1, nil
		}
		return 1, nil
	}
	if cp != rp {
		if cp < rp {
			return -1, nil
		}
		return 1, nil
	}
	return 0, nil
}
