package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var stdin = make(chan string, 1)

func main() {
	defer Perf()()
	/*
		Read stdin stream.
	*/
	go Scan()
	/*
		Validate the arguments to the program.
	*/
	if len(os.Args) < 2 {
		return
	}
	/*
		Parse the stop watch duration entered by the user.
	*/
	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
		Start the stop watch.
	*/
	deadLine := time.Now().UTC().Add(duration)
	go func() {
		for input := range stdin {
			if len(input) <= 0 {
				continue
			}
			switch input[0] {
			case '+':
				if len(input) < 2 {
					continue
				}
				duration, err := time.ParseDuration(input[1:])
				if err != nil {
					fmt.Println(err)
					continue
				}
				deadLine = deadLine.Add(duration)
			case '-':
				if len(input) < 2 {
					continue
				}
				duration, err := time.ParseDuration(input[1:])
				if err != nil {
					fmt.Println(err)
					continue
				}
				deadLine = deadLine.Add(-duration)
			case 'r':
				fmt.Println(time.Until(deadLine).String())
			case 'x':
				deadLine = time.Now().UTC()
			}
		}
	}()
	for time.Now().UTC().Before(deadLine) {
	}
	fmt.Println("Time-up!")
	_, err = Exec("notify-send", "-i", "bash", "-u", "critical", "pm", "Time-up!")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Exec(name string, args ...string) (stdout string, stderr error) {
	cmd := exec.Command(name, args...)
	bytes, stderr := cmd.Output()
	if stderr != nil {
		return
	}
	stdout = string(bytes)
	return
}

func Scan() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin <- scanner.Text()
	}
}

func Perf() func() {
	start := time.Now().UTC()
	return func() {
		fmt.Println(time.Since(start))
	}
}
