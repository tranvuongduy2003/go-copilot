package main

import (
	"fmt"
	"os"
)

// Build information (injected via ldflags)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	fmt.Printf("Go Copilot API Server\n")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Println()
	fmt.Println("TODO: Implement server startup")
	os.Exit(0)
}
