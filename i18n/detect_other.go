//go:build !windows

package i18n

func detectWindowsLanguage() Language {
	return ""
}
