package cheflinuxoperations

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func isPackageInstalled(packageName string) (bool, error) {
	cmd := exec.Command("apt", "-qq", "list", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(output), fmt.Sprintf("ii  %s", packageName)), nil
}

func InstallPackage(packageName string) {
	installed, err := isPackageInstalled(packageName)
	if err != nil {
		log.Fatalf("Error checking package installation: %v\n", err)
		return
	}

	if installed {
		fmt.Printf("Package %s is already installed.\n", packageName)
	} else {
		// Run the apt-get install command
		cmd := exec.Command("sudo", "apt-get", "install", "-y", packageName)

		// Capture the output of the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to install package %s: %v\n", packageName, err)
		} else {
			fmt.Printf("Package %s installed successfully.\n", packageName)
			fmt.Println("output", string(output))
		}

	}
}

func InstallMultiPackage(pkgNames []string) {
	for _, pkg := range pkgNames {
		InstallPackage(pkg)
	}
}

func RemovePackage(packageName string) {
	installed, err := isPackageInstalled(packageName)
	if err != nil {
		log.Printf("Error checking package installation: %v\n", err)
		return
	}

	if installed {
		// Run the apt-get remove command
		cmd := exec.Command("sudo", "apt-get", "remove", "-y", packageName)

		// Capture the output of the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to remove package %s: %v\n", packageName, err)
		} else {
			fmt.Printf("Package %s removed successfully.\n", packageName)
			fmt.Println("output", string(output))
		}
	} else {
		fmt.Printf("Package %s is not installed.\n", packageName)
	}
}

func RemoveMultiPackage(pkgNames []string) {
	for _, pkg := range pkgNames {
		RemovePackage(pkg)
	}
}
