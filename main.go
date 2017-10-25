package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	cmd := "ps -eo pid,etime,command | grep '[s]leep' | awk -v OFS='\t' '$1=$1'"
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

			min, _ := strconv.Atoi(times[len(times)-1])
			sec, _ := strconv.Atoi(times[len(times)-2])
			if min*60+sec < 30 {
				break
			}

			killCmd := fmt.Sprintf("sudo kill -9 %s", pid)
			exec.Command("sh", "-c", killCmd).Run()
		}
		time.Sleep(time.Second * 10)
	}
}
