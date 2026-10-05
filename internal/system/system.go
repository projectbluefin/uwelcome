package system

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"

	i "uwelcome/internal"
)

type ImageInfo struct {
	ImageName   string `json:"image-name"`
	ImageRef    string `json:"image-ref"`
	ImageFlavor string `json:"image-flavor"`
	ImageVendor string `json:"image-vendor"`
	ImageTag    string `json:"image-tag"`
}

const infoFile = "/usr/share/ublue-os/image-info.json"

var OSName = getOSName()

func GetDesktop() string {
	desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	if desktop == "" {
		return "Unknown"
	}
	return desktop
}

// GetGreenbootInfo is a command to greenboot status
func GetGreenbootInfo() string {
	if exists, err := i.DoesFileExist("/etc/motd.d/boot-status"); !exists || err != nil {
		return ""
	}

	re := regexp.MustCompile(`status is GREEN`)

	isGreen := re.FindString("status is GREEN")
	if isGreen != "" {
		return "healthy"
	}

	cmd := exec.Command("cat", "/etc/motd.d/boot-status")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return "`" + string(output) + "`"

}

// func IsSignedImage(info ImageInfo) bool {
// 	return strings.Contains(info.ImageRef, "ostree-image-signed:docker://")
// }

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

// getOSName gets the OS name from /etc/os-release
func getOSName() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`NAME="(.*)"`)
	match := re.FindStringSubmatch(string(data))
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
