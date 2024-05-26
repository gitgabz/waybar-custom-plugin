package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	waybarSignal := flag.Bool("s", false, "Run once and trigger update through SIGRTMIN+N in waybar config")
	continuous := flag.Bool("continuous", true, "Update output in a continuous manner if no signal is defined (default behaviour)")
	flag.Parse()

	_ = continuous // Keep continuous described in the flag output, it is more a placeholder since we can not affect it.

	if *waybarSignal {

		// Waybar will restart the process by SIGRTMIN+N and or interval as configured
		triggerBySignal()

	} else {

		// run with constant output
		go loopOutput()

		// Wait to be closed by SIGINT or SIGTERM to close gracefully
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

		// Do cleanup after exit signal is caught
		<-exit
		fmt.Println("Exiting...")

	}

}

func triggerBySignal() {
	/*
		# Refresh by interrupt signal, restarts the process that has exited.
		# The process runs once, outputs and dies. Then it is triggered by interval through signal or directly by SIGRTMIN+8 signal "pkill -RTMIN+8 waybar"

		"custom/lunch": {
			"format": "{}",
				"max-length": 40,
				"tooltip": false,
				"signal": 8,
				"interval": 12,
				"exec": "~/.local/bin/waybar_plugin -s",
				"on-click": "pkill -RTMIN+8 waybar",
				"return-type": "json"
		}
	*/

	outputToWaybar()
}

func loopOutput() {
	/*
		# Refresh continuously by Newline json and repeat output at will.

			"custom/lunch": {
				"format": "{}",
				"max-length": 40,
				"tooltip": false,
				"exec": "~/.local/bin/waybar_plugin",
				"return-type": "json"
			}
	*/

	for {
		outputToWaybar()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := r.Intn(10)
		time.Sleep(time.Duration(n) * time.Second)
	}
}

func outputToWaybar() {
	wbr := waybarResponse{
		Text:       strconv.Itoa(rand.Int()),
		Tooltip:    "",
		Class:      "",
		Percentage: 0,
	}

	j, err := wbr.jsonOutput()
	if err != nil {
		log.Fatalf("Failed to create json message output")
	}

	fmt.Printf("%s\n", j)
}

func (w *waybarResponse) jsonOutput() (jb []byte, err error) {
	jb, err = json.Marshal(w)
	return
}

type waybarResponse struct {
	Text       string `json:"text"`
	Tooltip    string `json:"tooltip"`
	Class      string `json:"class"`
	Percentage int64  `json:"percentage"`
}
