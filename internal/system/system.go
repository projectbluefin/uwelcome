package system

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
)

type ImageInfo struct {
	ImageName   string `json:"image-name"`
	ImageRef    string `json:"image-ref"`
	ImageFlavor string `json:"image-flavor"`
	ImageVendor string `json:"image-vendor"`
	ImageTag    string `json:"image-tag"`
}

const infoFile = "/usr/share/ublue-os/image-info.json"
const HealthyState = "healthy"
const bootStatusFile = "/etc/motd.d/boot-status"

// GetGreenbootInfo is a command to greenboot status
func GetGreenbootInfo() string {
	data, err := os.ReadFile(bootStatusFile)
	if err != nil {
		return ""
	}
	text := strings.TrimSpace(string(data))
	if strings.Contains(text, "status is GREEN") {
		return HealthyState
	}
	return text
}

// GetImageInfo is a Universal Blue focused command that retrieves the system image reference from their image-info file
func GetImageInfo() ImageInfo {

	defaultInfo := ImageInfo{}

	data, err := os.ReadFile(infoFile)
	if err != nil {
		return defaultInfo
	}

	var info ImageInfo
	err = json.Unmarshal(data, &info)
	if err != nil {
		return defaultInfo
	}

	return info
}

// GetOSName gets the OS name from os-release
func GetOSName() string {
	possiblePaths := []string{
		"/etc/os-release",
		"/usr/lib/os-release",
	}

	var data []byte

	for _, p := range possiblePaths {
		if d, err := os.ReadFile(p); err == nil {
			data = d
			break
		}
	}
	if len(data) == 0 {
		return ""
	}
	re := regexp.MustCompile(`(?m)^NAME="?([^"\n]*)"?$`)
	match := re.FindStringSubmatch(string(data))
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
