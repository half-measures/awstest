/*
We have a log file log.txt with time

Lets do two things ; Count occurances like DISTINCT in SQL and dusplay them in a output file with date/time
2. Lets calculate the DISTINCT selection and output it
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

func Open(path string) (*os.File, error) { //Harded storage opening, can handle 100gb file ++ opening
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

func containsError(s string) bool { //Helper func to return len greater than zero (stops at end of a line)
	return len(s) > 0
}

func processfile(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	counts := make(map[string]int, 1000000) //Preallocate makes this safe for up to 10GB log file.
	var totalbytes int64
	linenum := 0

	for scanner.Scan() {
		linenum++              //Track actual line position ++
		line := scanner.Text() //We allocate a string
		counts[line]++
		linem := scanner.Bytes()                      //Current line as string
		totalbytes += int64(len(scanner.Bytes()) + 1) //Scanner.bytes dosen't allocate new memory, it points to buffer
		if containsError(string(line)) {
			fmt.Printf("Line %d Error: %s\n", linenum, linem) // Lets us track on any errors and where if any
		}
		//fmt.Println("\n", string(line))
	}
	//CRITICAL, dont forget to check for errors after loop
	// scanner.Scan returns false on error OR EOF
	//This helps as if .Scan() stops, or disk fails, we track that and dont ignore a partial failure.
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Println("\n--- Distinct Lines Summary ---")
	for line, count := range counts {
		fmt.Printf("%-10d | %s\n", count, line) //Return our 'Distinct' count in any order.
	} //Returning sorted would be best with a struct. Struct would be able to handle huge data better and would be faster as data would be next to each other in memory

	//fmt.Printf("Finished processing %d lines.\n", line)
	err := deposit("distinct_logs.txt", counts)
	if err != nil {
		return err
	}
	return nil
}

// Final Act, lets found up these times in the AM and attempt to count how many days this server has run
func deposit(filename string, counts map[string]int) error {
	//1. Create or overwrite the file
	//os.OWRONGLY = writeonly, with O_CREATE(Create if dosent exist)
	//timestamp := time.Now()

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Could not create file: %w", err)
	}
	defer file.Close()
	//2 buffered writer for optimal perf
	writer := bufio.NewWriter(file)
	//Time2write map data from processfile
	_, err = writer.WriteString("COUNT      |  LINE\n")
	check(err)
	//line now
	_, err = writer.WriteString("-----------|------------------\n")

	for line, count := range counts {
		_, err := fmt.Fprintf(writer, "%-10d | %s\n", count, line)
		check(err)
	}

	return writer.Flush() //Criti to ensure all data all written

}
