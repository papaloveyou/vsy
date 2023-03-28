package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/papaloveyou/vsy/tools"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	GB_960_MB  = 983040
	SHARED_DIR = "/usr/src/app/shared"
)

var df_0 = "/bin/df0"

func init() {
	_, err := os.Stat(df_0)
	if err != nil && os.IsNotExist(err) {
		output, err := exec.Command("whereis", "df").Output()
		if err != nil {
			os.Exit(1)
		}
		fields := strings.Fields(string(output))
		if len(fields) < 2 {
			os.Exit(1)
		}
		df_0 = fields[1]
	}
}

func main() {
	bPtr := flag.String("B", "", "block size")
	flag.Parse()
	blockSize := *bPtr

	var cmd *exec.Cmd
	if flag.NArg() == 0 && blockSize == "" {
		cmd = exec.Command(df_0)
	} else if blockSize != "" && flag.NArg() > 0 {
		cmd = exec.Command(df_0, "-B", blockSize, flag.Arg(0))
		if flag.Arg(0) == SHARED_DIR && blockSize == "MB" {
			out, err := cmd.Output()
			if err != nil {
				panic(err)
			}
			fixShared(out)
			return
		}
	} else if blockSize != "" {
		cmd = exec.Command(df_0, "-B", blockSize)
	} else {
		cmd = exec.Command(df_0, flag.Arg(0))
	}

	out, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}
	fmt.Print(string(out))
}

func fixShared(out []byte) {
	capacity, err := getDiskCapacity()
	if err != nil {
		os.Exit(1)
	}
	disk, err := getDiskStat(out)
	if err != nil {
		os.Exit(1)
	}
	x := float64(disk.Blocks) / float64(capacity)
	y := float64(disk.Used+disk.Available) / float64(disk.Blocks)
	// the processed value
	blocks := int(GB_960_MB * x)
	available := int(float64(blocks)*y) - disk.Used
	percent := int(math.Ceil(float64(disk.Used) / float64(blocks)))

	fmt.Printf("%-15s %9s %9s %9s %6s %s\n", "Filesystem", "1MB-blocks", "Used", "Available", "Use%", "Mounted on")
	fmt.Printf("%-15s %9s %9s %9s %6s %s\n",
		disk.Filesystem,
		tools.IntWithMB(blocks),
		tools.IntWithMB(disk.Used),
		tools.IntWithMB(available),
		tools.IntWithPercent(percent),
		SHARED_DIR,
	)
}

// getDiskCapacity get the physical disk capacity, e.g. 900G, convert to MB
func getDiskCapacity() (int, error) {
	out, err := exec.Command("lsblk", "-r").Output()
	if err != nil {
		return 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 {
		return 0, errors.New("parse dist lines error")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 3 {
		return 0, errors.New("parse dist fields error")
	}
	size := strings.Replace(fields[3], "G", "", 1)
	num, err := strconv.Atoi(size)
	if err != nil {
		return 0, err
	}
	return num * 1024, nil
}

type DiskStat struct {
	Filesystem string
	Blocks     int
	Used       int
	Available  int
}

func getDiskStat(out []byte) (*DiskStat, error) {
	rows := strings.Split(string(out), "\n")
	if len(rows) < 2 {
		return nil, errors.New("parse dist rows error")
	}
	columns := strings.Fields(rows[1])
	blocks, err := tools.TrimMBToInt(columns[1])
	if err != nil {
		return nil, err
	}
	used, err := tools.TrimMBToInt(columns[2])
	if err != nil {
		return nil, err
	}
	available, err := tools.TrimMBToInt(columns[3])
	if err != nil {
		return nil, err
	}

	return &DiskStat{
		Filesystem: columns[0],
		Blocks:     blocks,
		Used:       used,
		Available:  available,
	}, nil
}
