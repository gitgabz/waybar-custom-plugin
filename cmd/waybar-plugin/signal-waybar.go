package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
)

const SIGRTMIN = 34

/*
	// Sample to call waybar process and restart by SIGRTMIN+N

	waybarSignal := flag.Int("s", 0, "Trigger with signal to waybar, the number is valid between 1 and N")

	pid, err := pidOfCommand("waybar")
	if err != nil {
		log.Fatalf("Failed to find PID of waybar: %s", err.Error())
	}

	sendSig(pid, waybarSignal)
*/

func sendSig(pid int, waybarSignal *int) {

	// Define the signal as SIGRTMIN + N
	sig := SIGRTMIN + *waybarSignal

	// Send signal to the process
	err := syscall.Kill(pid, syscall.Signal(sig))
	if err != nil {
		fmt.Printf("Failed to send signal: %s\n", err)
		os.Exit(1)
	}

}

// pidOfCommand loops through all commands and matches "command {-option1 -option2}", but not "command_extra, command-extra" to not confuse the pid we have found with a similarly named cmdline.
func pidOfCommand(nameOfCommand string) (pid int, err error) {

	// Define the regular expression pattern
	pattern := fmt.Sprintf(`^%s\b(?:\s+[a-z-A-Z0-9]\w+)*\s*$`, nameOfCommand)

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regular expression:", err)
		return
	}

	pid, err = findPIDByRegexp(re)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

// getCmdline reads the command line of a process from /proc/[pid]/cmdline
func getCmdline(pid string) (string, error) {
	cmdlinePath := filepath.Join("/proc", pid, "cmdline")
	cmdline, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return "", err
	}
	// Replace null bytes from /proc/{pid}/cmdline with spaces
	return string(bytes.ReplaceAll(cmdline, []byte{0}, []byte(" "))), nil
}

// findPIDByRegexp searches for the PID of the process with the given name
func findPIDByRegexp(re *regexp.Regexp) (int, error) {
	procPath := "/proc"
	entries, err := os.ReadDir(procPath)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if entry.IsDir() && isNumeric(entry.Name()) {
			pid := entry.Name()
			cmdline, err := getCmdline(pid)
			if err != nil {
				continue
			}

			if re.MatchString(cmdline) {
				return strconv.Atoi(pid)
			}
		}
	}

	return 0, fmt.Errorf("process %s not found", "waybar")
}

// isNumeric checks if a string consists only of digits
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
