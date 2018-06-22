package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	counter := false
	for true {
		tempList := []string{}

		time.Sleep(50000000)
		output, err1 := exec.Command("ipfs", "refs", "local").Output()
		if err1 != nil {
			os.Stderr.WriteString(err1.Error())
		}
		s := string(output)

		hashes := strings.Split(s, "\n")
		list := []string{}

		if counter == true {
			d := difference(list, tempList)
			fmt.Println("Nothing Found")
			for i := range d {
				output, err1 := exec.Command("ipfs", "get", d[i]).Output()
				fmt.Println(output)
				if err1 != nil {
					os.Stderr.WriteString(err1.Error())
				}
			}
		} else {
			for i := range hashes {
				fmt.Println(hashes[i])
				list = append(list, hashes[i])
				output, err1 := exec.Command("ipfs", "get", hashes[i]).Output()
				fmt.Println(string(output))
				if err1 != nil {
					os.Stderr.WriteString(err1.Error())
				}
			}
		}
		counter = true
		tempList = list
	}
}

func difference(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}
