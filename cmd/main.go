package main

import (
	"fmt"
	"os"

	"github.com/blaze-d83/hardware-monitor-go/internal/server"
)

func main() {
	fmt.Println("Starting system monitor..")

	srv := server.NewServer()
	srv.MonitorHardware()

	if err := srv.Start(); err != nil {
		fmt.Println("Server failed to start:", err)
		os.Exit(1)
	}

}
