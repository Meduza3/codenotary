package deps

import (
	"codenotary/internal/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)


func (c *Client) GetLatestVersionByProjectId(projectID string) (string, error) {
	
	pv, err := c.GetPackage(projectID)
	if err != nil {
		return "", fmt.Errorf("Failed to GetPackage: %w", err)
	}

	defaults := make([]string, 0)
	for _, v := range pv.Versions {
		if v.IsDefault {
			defaults = append(defaults, v.VersionKey.Version)
		}
	}
	switch len(defaults) {
	case 1:
		return defaults[0], nil
	default:
		return latestSemVer(pv.Versions)
	}
}

func latestSemVer(versions []models.Version) (string, error) {
	var max string
	var max1 int
	var max2 int
	var max3 int

	for _, v := range versions {
		stripped := filterNumbersAndDots(v.VersionKey.Name)
		parts := strings.Split(stripped, ".")
		if len(parts) != 3 {
			continue 
		}

		major, err1 := strconv.Atoi(parts[0])
		minor, err2 := strconv.Atoi(parts[1])
		patch, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			continue 
		}

		
		if major > max1 || (major == max1 && minor > max2) || (major == max1 && minor == max2 && patch > max3) {
			max = v.VersionKey.Name
			max1, max2, max3 = major, minor, patch
		}

	}

	if max == "" {
		return "", fmt.Errorf("no valid semantic versions found")
	}
	return max, nil
}

func filterNumbersAndDots(input string) string {
	
	re := regexp.MustCompile(`[^\d.]`)
	
	return re.ReplaceAllString(input, "")
}
