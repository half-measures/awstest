package main

// Pure day counter based on the problem
import (
	"fmt"
	"math"
)

const dailycap = 10 //Server can only process '10' per day

func main() {
	loglines := 34
	result := calculateprocessingdays(loglines)
	fmt.Printf("Processing %d logs will take %d days.\n", loglines, result)

}

func calculateprocessingdays(logcount int) int {
	if logcount == 0 {
		return 0
	}
	//float64 ensures math.Ceil can be used and allow remainder
	days := math.Ceil(float64(logcount)) / float64(dailycap)
	return int(days)
}
