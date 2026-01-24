package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"log"
)

// Device represents a drive found by smartctl
type Device struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Type     string `json:"type"`
}

// ScanResult matches the JSON output of "smartctl --scan --json"
type ScanResult struct {
	Devices []Device `json:"devices"`
}

func main() {
	fmt.Println("Starting Vigil Agent...")
	fmt.Println("1. Scanning for devices...")

	// Run "smartctl --scan --json"
	cmd := exec.Command("smartctl", "--scan", "--json")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running smartctl scan: %v. (Do you have smartmontools installed?)", err)
	}

	// Parse the JSON list of devices
	var scan ScanResult
	if err := json.Unmarshal(output, &scan); err != nil {
		log.Fatalf("Error parsing scan json: %v", err)
	}

	fmt.Printf("   Found %d devices.\n", len(scan.Devices))

	// Loop through each device and get its health
	for _, dev := range scan.Devices {
		fmt.Printf("2. Checking health for %s (%s)...\n", dev.Name, dev.Type)
		
		// Run "smartctl -x --json /dev/sdX"
		healthCmd := exec.Command("smartctl", "-x", "--json", "--device", dev.Type, dev.Name)
		healthOut, err := healthCmd.Output()
		if err != nil {
			log.Printf("   Failed to check %s: %v\n", dev.Name, err)
			continue
		}

		fmt.Printf("   -> Success! Retrieved %d bytes of SMART data.\n", len(healthOut))
	}
	
	fmt.Println("Done.")
}