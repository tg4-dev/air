package nodeinfo

import (
	"os/exec"
	"runtime"
	"strings"
)

// TODO: come up with additional fields for GPUInfo
type GPUInfo struct {
	Name string `json:"name"`
}


func DetectGPUName() string {
	switch runtime.GOOS {
	case "linux":
		return gpuNameLinux()
	case "windows":
		return gpuNameWindows()
	case "darwin":
		return gpuNameDarwin()
	default:
		return "unknown"
	}
}

func gpuNameLinux() string {
	out, err := exec.Command("sh", "-c", "lspci | grep -i vga").Output()
	if err != nil {
		return "unknown"
	}
	line := strings.TrimSpace(string(out))
	if idx := strings.Index(line, ": "); idx != -1 {
		return strings.TrimSpace(line[idx+2:])
	}
	return line
}

func gpuNameWindows() string {
	out, err := exec.Command("powershell", "-Command",
		"(Get-CimInstance Win32_VideoController).Name").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func gpuNameDarwin() string {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "Chipset Model:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Chipset Model:"))
		}
	}
	return "unknown"
}
