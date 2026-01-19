/*
We have a log file log.txt with time

Lets do two things ; Count occurances like DISTINCT in SQL and dusplay them in a output file with date/time
2. Lets calculate the
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {

	//ctx := context.Background() Not needed right now,but could provide context in future.
	// Open the log file
	path := "log.txt"
	file, err := Open(path)
	if err != nil {
		fmt.Println("Error opening file, Please check file:", err)
		return

	}
	defer file.Close() // Caller defer outside the open to ensure cleanup
	fmt.Printf("File at %s opened successfully\n", path)
	err = processfile(file) //Note the cool use here of err, which allows a if check anyways!
	if err != nil {
		fmt.Println("Error Processing file, Please check file for corrption or malforms:", err)
	}

}

func Open(path string) (*os.File, error) {
	//Path sanitization can be added here
	cleanpath := filepath.Clean(path) //Clean the path to prevent directory traversal attacks

	//Lets check if file exists

	file, err := os.Open(cleanpath) //os.Open Opens as read only just in case
	if err != nil {
		//Below we handle specific errors
		if errors.Is(err, os.ErrNotExist) {
			// File does not exist, deliver errors saying that it may be elsewhere
			fmt.Errorf("Storage: Filepath %s does not exist: %w\n", cleanpath, err)
		}
		//Permission or sys errs
		if errors.Is(err, os.ErrPermission) {
			return nil, fmt.Errorf("Storage: Permission denied for %s: %w", cleanpath, err)
		}

		return nil, fmt.Errorf("Storage: Error opening file at %s: %v", cleanpath, err)

	}
	return file, nil
}

func check(e error) { //Helper Checker function
	if e != nil {
		panic(e)
	}
}

func containsError(s string) bool {
	return len(s) > 0
}

func processfile(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	linecount := 0
	var totalbytes int64

	for scanner.Scan() {
		linecount++
		line := scanner.Bytes() //Current line as string
		totalbytes += int64(len(line) + 1)
		if containsError(string(line)) {
			fmt.Printf("Line %d: %s\n", linecount, line) // Lets us track on any errors and where if any
		}
		//fmt.Println("\n", string(line))
	}
	//CRITICAL, dont forget to check for errors after loop
	// scanner.Scan returns false on error OR EOF
	//This helps as if .Scan() stops, or disk fails, we track that and dont ignore a partial failure.
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}

	fmt.Printf("Finished processing %d lines.\n", linecount)

	return nil
}

//Final Act, lets found up these times in the AM and attempt to count how many days this server has run
