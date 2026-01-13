// cmd/rdxbus/easy_prompt.go
package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func prompt(reader *bufio.Reader, label, def string) string {
	fmt.Printf("%s [%s]: ", label, def)
	in, _ := reader.ReadString('\n')
	in = strings.TrimSpace(in)
	if in == "" {
		return def
	}
	return in
}

func promptInt(reader *bufio.Reader, label string, def int) int {
	for {
		v := prompt(reader, label, strconv.Itoa(def))
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
		fmt.Println("Invalid number")
	}
}
