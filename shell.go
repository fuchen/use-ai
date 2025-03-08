package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// SystemInfo 存储系统和shell信息
type SystemInfo struct {
	OS        string
	ShellType string
	ShellPath string
}

// DetectSystem 检测当前的操作系统和shell类型
func DetectSystem() SystemInfo {
	info := SystemInfo{
		OS: runtime.GOOS,
	}

	// 检测shell类型
	switch runtime.GOOS {
	case "windows":
		// 先检查PowerShell特有的环境变量PSModulePath
		if os.Getenv("PSModulePath") != "" {
			info.ShellType = "powershell"
			// 尝试找到PowerShell路径
			possiblePaths := []string{
				"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
				"C:\\Program Files\\PowerShell\\7\\pwsh.exe",  // PowerShell Core路径
				os.ExpandEnv("${ProgramFiles}\\PowerShell\\7\\pwsh.exe"),
			}

			// 检查环境变量POWERSCRIPT_EXE
			if posh := os.Getenv("POWERSCRIPT_EXE"); posh != "" {
				possiblePaths = append([]string{posh}, possiblePaths...)
			}

			// 尝试找到可用的PowerShell路径
			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					info.ShellPath = path
					break
				}
			}

			// 如果没找到，使用命令获取
			if info.ShellPath == "" {
				cmd := exec.Command("where", "powershell")
				output, err := cmd.Output()
				if err == nil && len(output) > 0 {
					info.ShellPath = strings.TrimSpace(string(output))
				} else {
					// 如果还是找不到，使用默认名称
					info.ShellPath = "powershell.exe"
				}
			}
		} else {
			// 尝试获取当前shell环境
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = os.Getenv("ComSpec")
			}

			if strings.Contains(strings.ToLower(shell), "cmd") {
				info.ShellType = "cmd"
			} else {
				// 尝试检测git bash或其他类型的shell
				cmd := exec.Command("sh", "-c", "echo $SHELL")
				output, err := cmd.CombinedOutput()
				if err == nil && len(output) > 0 {
					outputStr := strings.TrimSpace(string(output))
					if strings.Contains(outputStr, "bash") {
						info.ShellType = "bash"
					} else {
						info.ShellType = "unknown"
					}
				} else {
					info.ShellType = "cmd" // 默认使用cmd
				}
			}
			info.ShellPath = shell
		}

	case "darwin", "linux":
		// 在Unix系统上，我们可以检查$SHELL环境变量
		shell := os.Getenv("SHELL")
		info.ShellPath = shell

		if strings.Contains(shell, "bash") {
			info.ShellType = "bash"
		} else if strings.Contains(shell, "zsh") {
			info.ShellType = "zsh"
		} else if strings.Contains(shell, "fish") {
			info.ShellType = "fish"
		} else {
			info.ShellType = "sh" // 默认使用sh
		}

	default:
		info.ShellType = "unknown"
		info.ShellPath = "unknown"
	}

	return info
}

// GetSystemDescription 返回系统和shell的描述信息
func (si SystemInfo) GetSystemDescription() string {
	return fmt.Sprintf("OS: %s, Shell Type: %s, Shell Path: %s", si.OS, si.ShellType, si.ShellPath)
}