package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	commandPtr := flag.String("cmd", "", "command executing in the process to kill")
	killSecPtr := flag.Int64("sec", 30, "process which exceeds this second will be killed")
	intervalPtr := flag.Int64("interval", 10, "interval of checking status and killing")
	flag.Parse()

	fmt.Println(*commandPtr, *killSecPtr, *intervalPtr)
	cmd := fmt.Sprintf("ps -eo pid,etime,command | grep '%s' | grep -v grep | awk -v OFS='\t' '$1=$1'", (*commandPtr))
	// cmd := "ps -eo pid,etime,command | grep '[s]leep' | awk -v OFS='\t' '$1=$1'"
	for {
		out, _ := exec.Command("sh", "-c", cmd).Output()

		lines := strings.Split(string(out), "\n")
		for _, line := range lines[:len(lines)-1] {
			arr := strings.Split(line, "\t")
			if len(arr) < 2 {
				break
			}

			pid := arr[0]
			etime := arr[1]
			times := strings.Split(etime, ":")
			if len(times) < 2 {
				break
			}

			min, _ := strconv.ParseInt(times[len(times)-1], 10, 64)
			sec, _ := strconv.ParseInt(times[len(times)-2], 10, 64)
			if min*60+sec < *killSecPtr {
				break
			}

			killCmd := fmt.Sprintf("sudo kill -9 %s", pid)
			exec.Command("sh", "-c", killCmd).Run()
		}
		time.Sleep(time.Second * time.Duration(*intervalPtr))
	}
}
