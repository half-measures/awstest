package main

// Modified program that only counts [ERROR] bits
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	filepath := "log.txt"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Unable to open file. Check if file exists or is malformed:%v", err)
	}
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		if strings.Contains(line, "[ERROR]") {
			count++
			fmt.Println(line)
		}

	}

}
