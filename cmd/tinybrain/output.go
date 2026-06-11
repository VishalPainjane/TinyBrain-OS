package main

import (
	"fmt"
	"io"
)

type checkLevel int

const (
	checkOK checkLevel = iota
	checkWarn
	checkFail
)

type checkResult struct {
	level   checkLevel
	message string
}

func printCheck(w io.Writer, level checkLevel, format string, args ...any) checkResult {
	prefix := "[ok]  "
	switch level {
	case checkWarn:
		prefix = "[warn]"
	case checkFail:
		prefix = "[fail]"
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(w, "  %s %s\n", prefix, msg)
	return checkResult{level: level, message: msg}
}

func worstLevel(results []checkResult) checkLevel {
	worst := checkOK
	for _, r := range results {
		if r.level > worst {
			worst = r.level
		}
	}
	return worst
}

func exitCodeForChecks(results []checkResult) int {
	switch worstLevel(results) {
	case checkFail:
		return 1
	default:
		return 0
	}
}
