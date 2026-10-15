package utils

import "fmt"

// FormatBytes format the bytes size int to human readable string
// 500 -> 500 B
// 150000  -> 146.5 KiB
// 1073741824 -> 1.0 GiB
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
