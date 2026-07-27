package scriptmgr

import (
	"embed"
	"fmt"
)

//go:embed all:scripts
var scriptFiles embed.FS

// GetScriptBytes 直接按需从内嵌资源中读取字节，零多余状态
func GetScriptBytes(scriptName string) ([]byte, error) {
	filePath := fmt.Sprintf("scripts/%s", scriptName)
	data, err := scriptFiles.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取内嵌脚本失败 [%s]: %v", scriptName, err)
	}
	return data, nil
}
