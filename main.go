package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/reddydinesh427/upgraded-config-mang/pkg/linuxoperations"
	"golang.org/x/crypto/ssh"
)

type InventoryDetails struct {
	Inventory struct {
		Hosts    []string `json:"hosts"`
		Username string   `json:"username"`
		Password string   `json:"password"`
	} `json:"inventory"`
}

type Task struct {
	Action         string   `json:"action"`
	PackageName    []string `json:"packageNames"`
	Path           string   `json:"path"`
	Content        string   `json:"content"`
	LocalFilePath  string   `json:"localPath"`
	RemoteFilePath string   `json:"remotePath"`
	Owner          string   `json:"owner"`
	Group          string   `json:"group"`
	Perm           string   `json:"perm"`
	DirectoryPath  string   `json:"directoryPath"`
	ServiceName    string   `json:"serviceName"`
	Depends        struct {
		ServiceName string `json:"serviceName"`
	} `json:"depends,omitempty"`
}
type TaskConfig struct {
	Tasks []Task `json:"tasks"`
}

var (
	InventoryFile = flag.String("inventory", "", "inventory file path, must be a JSON, EG: ./inventory.json")
	TasksFile     = flag.String("tasks", "", "Task file path, must be a JSON EG: ./tasks.json")
)

func init() {
	flag.Parse()
	if *InventoryFile == "" || *TasksFile == "" {
		fmt.Println("Inventory file and task file are missing")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	hostsData, err := ioutil.ReadFile(*InventoryFile)
	if err != nil {
		log.Fatal("Error reading hosts.json:", err)
	}
	var hosts InventoryDetails
	err = json.Unmarshal(hostsData, &hosts)
	if err != nil {
		log.Fatal("Error decoding hosts.json:", err)
	}
	tasksData, err := ioutil.ReadFile(*TasksFile)
	if err != nil {
		log.Fatal("Error reading tasks.json:", err)
	}
	var taskConfigs TaskConfig
	err = json.Unmarshal(tasksData, &taskConfigs)
	if err != nil {
		log.Fatal("Error decoding tasks.json:", err)
	}

	for _, host := range hosts.Inventory.Hosts {
		sshConfig := &ssh.ClientConfig{
			User: hosts.Inventory.Username,
			Auth: []ssh.AuthMethod{
				ssh.Password(hosts.Inventory.Password),
				// Add SSH key authentication method if preferred
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		remoteServer := host + ":22"

		conn, err := ssh.Dial("tcp", remoteServer, sshConfig)
		if err != nil {
			log.Fatalf("Failed to dial: %v", err)
		}
		defer conn.Close()

		for _, task := range taskConfigs.Tasks {

			// fmt.Println(host, task.Action)
			switch task.Action {
			case "packageInstall":
				//logic for instal
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				linuxoperations.InstallMultiPackage(conn, task.PackageName)
			case "packageRemove":
				//logic for instal
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				linuxoperations.RemoveMultiPackage(conn, task.PackageName)
			case "remoteFile":
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				diffChecker := linuxoperations.RemoteFileCopy(conn, task.LocalFilePath, task.RemoteFilePath, task.Owner, task.Group, task.Perm)
				if task.Depends.ServiceName != "" && diffChecker {
					fmt.Println("Restarting service:", task.Depends.ServiceName)
					linuxoperations.ServiceRestart(conn, task.Depends.ServiceName)
				}
			case "directory":
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				linuxoperations.CreateDirectory(conn, task.DirectoryPath, task.Perm, task.Owner, task.Group)
			case "serviceStart":
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				linuxoperations.ServiceStart(conn, task.ServiceName)
			case "serviceStop":
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				linuxoperations.ServiceStop(conn, task.ServiceName)
			case "serviceRestart":
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				linuxoperations.ServiceRestart(conn, task.ServiceName)
			case "removeFile":
				fmt.Printf("\nRunning on Host: %s, Action: %s\n", host, task.Action)
				linuxoperations.RemoveFile(conn, task.Path)
			default:
				fmt.Println("No action was defined")
			}
			// if task.Action == "install" {
			// 	// Start a session for each command
			// 	session, err := conn.NewSession()
			// 	var outputBuilder strings.Builder
			// 	outputWriter := io.MultiWriter(log.Writer(), &outputBuilder)

			// 	// Set up pipes for input/output
			// 	session.Stdout = outputWriter
			// 	session.Stderr = outputWriter
			// 	if err != nil {
			// 		log.Fatalf("Failed to create session: %v", err)
			// 	}
			// 	// Service management commands
			// 	serviceName := "cron"
			// 	// startServiceCmd := fmt.Sprintf("sudo systemctl start %s", serviceName)
			// 	startServiceCmd := fmt.Sprintln("sudo systemctl list-units|cat")
			// 	serviceStatusCmd := fmt.Sprintf("sudo systemctl is-active %s", serviceName)
			// 	// Check if service is active
			// 	err = session.Run(serviceStatusCmd)
			// 	if err != nil {
			// 		log.Printf("Error checking service status: %v", err)
			// 	} else {
			// 		capturedOutput := outputBuilder.String()
			// 		if strings.Contains(capturedOutput, "active") {
			// 			// Start the service
			// 			fmt.Println("Starting service:", serviceName)
			// 			err = session.Run(startServiceCmd)
			// 			if err != nil {
			// 				log.Printf("Failed to start service: %v", err)
			// 			} else {
			// 				fmt.Println("Service started.")
			// 			}
			// 		} else {
			// 			fmt.Println("Service is already active.")
			// 		}
			// 	}

			// }
		}
	}

}
