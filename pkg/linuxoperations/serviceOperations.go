package linuxoperations

import (
	"fmt"
	"strings"

	"github.com/reddydinesh427/upgraded-config-mang/pkg/helper"
	"golang.org/x/crypto/ssh"
)

func ServiceStart(conn *ssh.Client, serviceName string) {
	checkServiceStatusCmd := fmt.Sprintf("sudo systemctl is-active %s", serviceName)
	serviceStatus := helper.CmdRunner(conn, checkServiceStatusCmd)
	if !strings.Contains(serviceStatus, "active") {
		fmt.Println("Starting service:", serviceName)
		startServiceCmd := fmt.Sprintf("sudo systemctl start %s", serviceName)
		helper.CmdRunner(conn, startServiceCmd)
		fmt.Println("Service started.")
	} else {
		fmt.Println("Service is already active.")
	}
}

func ServiceStop(conn *ssh.Client, serviceName string) {
	checkServiceStatusCmd := fmt.Sprintf("sudo systemctl is-active %s", serviceName)
	serviceStatus := helper.CmdRunner(conn, checkServiceStatusCmd)
	if strings.Contains(serviceStatus, "active") {
		fmt.Println("Starting service:", serviceName)
		startServiceCmd := fmt.Sprintf("sudo systemctl stop %s", serviceName)
		helper.CmdRunner(conn, startServiceCmd)
		fmt.Println("Service started.")
	} else {
		fmt.Println("Service is inactive, no further action needed.")
	}
}

func ServiceRestart(conn *ssh.Client, serviceName string) {

	fmt.Println("Re-Starting service:", serviceName)
	startServiceCmd := fmt.Sprintf("sudo systemctl restart %s", serviceName)
	helper.CmdRunner(conn, startServiceCmd)
	fmt.Println("Service restarted.")

}
