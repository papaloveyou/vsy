package internal

import (
	"github.com/papaloveyou/vsy/tools"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	MIN_MEMORY_GB = 30
	GB_32_KB      = 33554432
)

var meminfoRegx = regexp.MustCompile(`\d{2,}`)

func GetMeminfo() string {
	info, err := readMeminfo()
	if err != nil {
		log.Println(err)
		return ""
	}
	lines := strings.Split(info, "\n")
	memTotal := getMemTotal(lines[0])
	log.Println("memTotal:", memTotal)

	if tools.KBToGB(memTotal) >= MIN_MEMORY_GB {
		return info
	}
	gb := tools.KbToRoundedUpGb(memTotal)
	log.Println("KbToRoundedUpGb:", gb)
	genMemTotal := int(GB_32_KB * tools.CalculateMultiple(memTotal, gb*1024*1024))
	log.Println("genMemTotal:", genMemTotal, tools.KBToGB(genMemTotal))

	multiple := tools.CalculateMultiple(genMemTotal, memTotal)
	log.Println("multiple:", multiple)

	for i, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			lines[i] = replaceLineValue(line, genMemTotal)
			continue
		}
		if strings.HasPrefix(line, "MemFree:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "Buffers:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "Cached:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "Active(file):") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "Inactive(file):") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "KReclaimable:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "SReclaimable:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "CommitLimit:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
		if strings.HasPrefix(line, "DirectMap2M:") {
			lines[i] = fixLineValue(line, multiple)
			continue
		}
	}
	return strings.Join(lines, "\n")
}

func fixLineValue(line string, multiple float64) string {
	value := parseLineValue(line)
	value = int(float64(value) * multiple)
	line = replaceLineValue(line, value)
	return line
}

func parseLineValue(line string) int {
	find := meminfoRegx.FindString(line)
	find = strings.TrimSpace(find)
	number, err := strconv.Atoi(find)
	if err != nil {
		return 0
	}
	return number
}

func replaceLineValue(line string, value int) string {
	return meminfoRegx.ReplaceAllString(line, strconv.Itoa(value))
}

func readMeminfo() (string, error) {
	bytes, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}
	data := string(bytes)
	return data, nil
}

func getMemTotal(line string) int {
	val := strings.Fields(line)[1]
	number, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return number
}
