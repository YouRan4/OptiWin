package main

import (
	"strconv"
	"strings"
)

const CurrentVersion = "v1.2.2"

func parseVersion(v string) (major, minor int) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) >= 1 {
		major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		minor, _ = strconv.Atoi(parts[1])
	}
	return
}

func compareVersion(current, remote string) int {
	cm, cn := parseVersion(current)
	rm, rn := parseVersion(remote)
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
	return 0
}
