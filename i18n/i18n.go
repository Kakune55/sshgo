package i18n

import (
	"fmt"
	"os"
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
	// 首先检查环境变量
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	
	// 如果环境变量包含中文相关设置，则使用中文
	if lang != "" && (strings.Contains(lang, "zh") || strings.Contains(lang, "ZH")) {
		return Chinese
	}
	
	// 检查系统区域设置
	// 在Windows系统中检查系统区域
	if os.Getenv("OS") == "Windows_NT" {
		// 简单检查，如果是在中文Windows系统上，默认使用中文
		// 更复杂的检测可以通过检查系统区域设置实现
		if strings.Contains(os.Getenv("NUMBER_OF_PROCESSORS"), "Chinese") ||
		   strings.Contains(os.Getenv("PROCESSOR_IDENTIFIER"), "Chinese") {
			return Chinese
		}
	}
	
	// 默认使用英文
	return English
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