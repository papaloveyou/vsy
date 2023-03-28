package main

import (
	"fmt"
	. "github.com/papaloveyou/vsy/internal"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var numberRegx = regexp.MustCompile(`\d`)
var stringRegx = regexp.MustCompile(`0-\d`)

func main() {
	cmd := exec.Command("lscpu")
	output, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}
	numCPU := runtime.NumCPU()
	if numCPU >= MIN_CPU_CORES {
		fmt.Println(string(output))
		return
	}
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "CPU(s):") {
			line = replaceNumberValue(line, MIN_CPU_CORES)
			lines[i] = line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "On-line CPU(s) list:") {
			line = replaceStringValue(line, "0-5")
			lines[i] = line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "Core(s) per socket:") ||
			strings.HasPrefix(strings.TrimSpace(line), "Core(s) per cluster:") {
			line = replaceNumberValue(line, MIN_CPU_CORES)
			lines[i] = line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "NUMA node0 CPU(s):") {
			line = replaceStringValue(line, "0-5")
			lines[i] = line
		}
	}
	fmt.Println(strings.Join(lines, "\n"))
}

func replaceNumberValue(line string, value int) string {
	return numberRegx.ReplaceAllString(line, strconv.Itoa(value))
}

func replaceStringValue(line string, value string) string {
	return stringRegx.ReplaceAllString(line, value)
}
