//go:build windows

package i18n

import "syscall"

// detectWindowsLanguage 通过 Windows API 获取系统 UI 语言
func detectWindowsLanguage() Language {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetUserDefaultUILanguage")
	langID, _, _ := proc.Call()

	// 取低 10 位为 Primary Language ID
	// LANG_CHINESE = 0x04
	if langID&0x3FF == 0x04 {
		return Chinese
	}

	return ""
}
