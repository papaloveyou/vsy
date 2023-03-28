package internal

import (
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	MIN_CPU_CORES = 6
)

var numberRegx = regexp.MustCompile(`\d`)

func GetCpuinfo() string {
	info := readCpuinfo()
	numCPU := runtime.NumCPU()
	log.Println("NumCPU:", numCPU)
	if numCPU >= MIN_CPU_CORES {
		return info
	}

	lines := strings.Split(info, "\n")
	var processor []string
	for _, line := range lines {
		if line == "" || strings.TrimSpace(line) == "" {
			processor = append(processor, line)
			break
		}
		if strings.HasPrefix(line, "siblings") {
			line = replaceNumberValue(line, MIN_CPU_CORES)
		}
		if strings.HasPrefix(line, "cpu cores") {
			line = replaceNumberValue(line, MIN_CPU_CORES)
		}
		processor = append(processor, line)
	}

	var processors []string
	for i := 0; i < MIN_CPU_CORES; i++ {
		for _, line := range processor {
			if strings.HasPrefix(line, "processor") {
				line = replaceNumberValue(line, i)
			}
			if strings.HasPrefix(line, "core id") {
				line = replaceNumberValue(line, i)
			}
			if strings.HasPrefix(line, "apicid") {
				line = replaceNumberValue(line, i)
			}
			if strings.HasPrefix(line, "initial apicid") {
				line = replaceNumberValue(line, i)
			}
			processors = append(processors, line)
		}
	}

	info = strings.Join(processors, "\n")
	return info
}

func readCpuinfo() string {
	bytes, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return ""
	}
	data := string(bytes)
	return data
}

func replaceNumberValue(line string, value int) string {
	return numberRegx.ReplaceAllString(line, strconv.Itoa(value))
}
