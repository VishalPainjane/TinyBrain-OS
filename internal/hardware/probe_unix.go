//go:build !windows

package hardware

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func detectRAMBytes() (uint64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("read /proc/meminfo: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				break
			}
			kb, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("parse MemTotal: %w", err)
			}
			return kb * 1024, nil
		}
	}

	return 0, fmt.Errorf("MemTotal not found in /proc/meminfo")
}
