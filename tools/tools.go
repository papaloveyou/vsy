package tools

import (
	"math"
	"os"
	"strconv"
	"strings"
)

func MkDir(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenFsname() string {
	const prefix = "/dev/sd"
	var suffixes = [26]string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}

	for _, suffix := range suffixes {
		fsname := prefix + suffix
		if _, err := os.Stat(fsname); err != nil && os.IsNotExist(err) {
			return fsname
		}
	}
	return "/dev/vda1"
}

func PString(v string) *string {
	return &v
}

// KBToGB if it's less than 1 GiB, return 1 GiB
func KBToGB(kb int) int {
	gb := kb / 1024 / 1024
	if gb > 0 {
		return gb
	}
	return 1
}

func KbToRoundedUpGb(kb int) int {
	gb := float64(kb) / (1024 * 1024)
	roundedUpGb := math.Ceil(gb)
	return int(roundedUpGb)
}

func CalculateMultiple(dividend, divisor int) float64 {
	return float64(dividend) / float64(divisor)
}

// TrimMBToInt 100MB -> 100
func TrimMBToInt(str string) (int, error) {
	return trimUnitTotoInt(str, "MB")
}

// IntWithMB 100 -> 100MB
func IntWithMB(i int) string {
	return strconv.Itoa(i) + "MB"
}

func IntWithPercent(i int) string {
	return strconv.Itoa(i) + "%"
}

func trimUnitTotoInt(str string, unit string) (int, error) {
	size := strings.Replace(str, unit, "", 1)
	num, err := strconv.Atoi(size)
	if err != nil {
		return 0, err
	}
	return num, nil
}
