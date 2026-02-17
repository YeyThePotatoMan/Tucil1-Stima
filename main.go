package main

import (
	"bufio"
	"fmt"
	"os"
)

func read_file(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}

	sc := bufio.NewScanner(f)
	g = []string{}
	for sc.Scan() {
		cnt += 1
		line := sc.Text()
		if len(line) > 0 {
			g = append(g, line)
		}
	}
	// TODO : add validation

	n = len(g)
	ans = make([]int, n)
	for i := range ans {
		ans[i] = -1
	}
	return true
}

func main() {
	var path string
	fmt.Scan(&path)

	read_file(path)
	// fmt.Scan(&n)

	// for i := 0; i < n; i++ {
	// 	fmt.Scan(&g[i])
	// }

	mask := make(map[uint8]bool)
	solve2(0, false, mask)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j == ans[i] {
				fmt.Print("#")
			} else {
				fmt.Print(string(g[i][j]))
			}
		}
		fmt.Println()
	}
}
