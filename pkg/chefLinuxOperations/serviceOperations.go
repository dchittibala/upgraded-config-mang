package cheflinuxoperations

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func isServiceRunning(serviceName string) (bool, error) {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(output)) == "active", nil
}

func StartService(serviceName string) {
	running, err := isServiceRunning(serviceName)
	if err != nil {
		log.Printf("Error checking service status: %v\n", err)
		return
	}

	if running {
		fmt.Printf("Service %s is already running.\n", serviceName)
	} else {
		// Run the systemctl start command
		cmd := exec.Command("sudo", "systemctl", "start", serviceName)

		// Capture the output of the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to start service %s: %v\n", serviceName, err)
		} else {
			fmt.Printf("Service %s started successfully.\n", serviceName)
			fmt.Println("output", output)
		}

	}
}

func StopService(serviceName string) {
	running, err := isServiceRunning(serviceName)
	if err != nil {
		log.Printf("Error checking service status: %v\n", err)
		return
	}

	if running {
		// Run the systemctl stop command
		cmd := exec.Command("sudo", "systemctl", "stop", serviceName)

		// Capture the output of the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to stop service %s: %v\n", serviceName, err)
		} else {
			fmt.Printf("Service %s stopped successfully.\n", serviceName)
			fmt.Println("output", output)
		}
	} else {
		fmt.Printf("Service %s is already stopped.\n", serviceName)
	}
}

func RestartService(serviceName string) {
	cmd := exec.Command("sudo", "systemctl", "restart", serviceName)

	// Capture the output of the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to stop service %s: %v\n", serviceName, err)
	} else {
		fmt.Printf("Service %s stopped successfully.\n", serviceName)
		fmt.Println("output", output)
	}
}
