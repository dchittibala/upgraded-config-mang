package linuxoperations

import (
	"fmt"
	"strings"

	"github.com/reddydinesh427/upgraded-config-mang/pkg/helper"
	"golang.org/x/crypto/ssh"
)

func InstallPackage(conn *ssh.Client, packageName string) {
	checkServiceCmd := fmt.Sprintf("apt -qq list %s", packageName)
	serviceStatus := helper.CmdRunner(conn, checkServiceCmd)
	if !strings.Contains(serviceStatus, "installed") {
		fmt.Println("Installing package:", packageName)
		startServiceCmd := fmt.Sprintf("sudo apt install -y %s", packageName)
		helper.CmdRunner(conn, startServiceCmd)
		fmt.Println("package installed.")
	} else {
		fmt.Println("Package already installed.")
	}
}

func InstallMultiPackage(conn *ssh.Client, packageName []string) {
	if len(packageName) == 0 {
		panic("No packages were found, Please check the name of the key for packageNames for action: packageInstall in tasks.json")
	}
	for _, pkg := range packageName {
		InstallPackage(conn, pkg)
	}
}

func RemovePackage(conn *ssh.Client, packageName string) {
	checkServiceCmd := fmt.Sprintf("apt -qq list %s", packageName)
	serviceStatus := helper.CmdRunner(conn, checkServiceCmd)
	if strings.Contains(serviceStatus, "installed") {
		fmt.Println("Installing package:", packageName)
		startServiceCmd := fmt.Sprintf("sudo apt remove -y %s", packageName)
		helper.CmdRunner(conn, startServiceCmd)
		fmt.Println("package removed.")
	} else {
		fmt.Println("Package doesnt exist.")
	}
}
func RemoveMultiPackage(conn *ssh.Client, packageName []string) {
	if len(packageName) == 0 {
		panic("No packages were found, Please check the name of the key for packageNames for action: packageRemove in tasks.json")
	}
	for _, pkg := range packageName {
		RemovePackage(conn, pkg)
	}
}
