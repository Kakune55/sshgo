package i18n

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Language 定义支持的语言
type Language string

const (
	Chinese Language = "zh"
	English Language = "en"
)

// 当前语言
var currentLanguage Language

// 初始化包时设置默认语言
func init() {
	currentLanguage = detectLanguage()
}

// detectLanguage 检测系统语言
func detectLanguage() Language {
	// 优先使用显式环境变量覆盖（更可控）
	if explicit := os.Getenv("SSHGO_LANG"); explicit != "" {
		switch strings.ToLower(explicit) {
		case "zh", "zh_cn", "zh-cn", "cn", "chinese":
			return Chinese
		case "en", "en_us", "en-us", "english":
			return English
		}
	}

	// 其次检查系统常见语言环境变量
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	l := strings.ToLower(lang)
	if l != "" && strings.Contains(l, "zh") {
		return Chinese
	}

	// Windows 获取语言设置
	if runtime.GOOS == "windows" {
		if wl := detectWindowsLanguage(); wl != "" {
			return wl
		}
	}
	return English
}

// SetLanguage 允许在运行时显式切换语言（如后续想做命令行参数 --lang=en）
func SetLanguage(lang Language) {
	switch lang {
	case Chinese, English:
		currentLanguage = lang
	}
}

// T 获取指定键的翻译字符串
func T(key StringKey) string {
	var translations map[StringKey]string

	switch currentLanguage {
	case Chinese:
		translations = zhStrings
	case English:
		translations = enStrings
	default:
		translations = zhStrings // 默认使用中文
	}

	if translation, exists := translations[key]; exists {
		return translation
	}

	// 如果找不到翻译，返回键名
	return string(key)
}

// TWithArgs 获取带参数的翻译字符串
func TWithArgs(key StringKey, args ...interface{}) string {
	format := T(key)
	return fmt.Sprintf(format, args...)
}
