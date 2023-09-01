package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cheflinuxoperations "github.com/reddydinesh427/upgraded-config-mang/pkg/chefLinuxOperations"
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
	TasksFile = flag.String("tasks", "", "Task file path, must be a JSON EG: ./tasks.json")
)

func init() {
	flag.Parse()
	if *TasksFile == "" {
		fmt.Println("task file are missing")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {

	tasksData, err := ioutil.ReadFile(*TasksFile)
	if err != nil {
		log.Fatal("Error reading tasks.json:", err)
	}
	var taskConfigs TaskConfig
	err = json.Unmarshal(tasksData, &taskConfigs)
	if err != nil {
		log.Fatal("Error decoding tasks.json:", err)
	}

	for _, task := range taskConfigs.Tasks {
		switch task.Action {
		case "packageInstall":
			//logic for instal
			fmt.Printf("\nRunning on  Action: %s\n", task.Action)
			cheflinuxoperations.InstallMultiPackage(task.PackageName)
		case "packageRemove":
			//logic for instal
			fmt.Printf("\nRunning on  Action: %s\n", task.Action)
			cheflinuxoperations.RemoveMultiPackage(task.PackageName)
		case "remoteFile":
			fmt.Printf("\nRunning on  Action: %s\n", task.Action)
			diffChecker := cheflinuxoperations.RemoteFile(task.LocalFilePath, task.RemoteFilePath, task.Owner, task.Group, task.Perm)
			if task.Depends.ServiceName != "" && diffChecker {
				fmt.Println("Restarting service:", task.Depends.ServiceName)
				cheflinuxoperations.RestartService(task.Depends.ServiceName)
			}
		case "serviceStart":
			fmt.Printf("\nRunning on  Action: %s\n", task.Action)
			cheflinuxoperations.StartService(task.ServiceName)
		case "serviceStop":
			fmt.Printf("\nRunning on  Action: %s\n", task.Action)
			cheflinuxoperations.StopService(task.ServiceName)
		case "serviceRestart":
			fmt.Printf("\nRunning on  Action: %s\n", task.Action)
			cheflinuxoperations.RestartService(task.ServiceName)
		case "removeFile":
			fmt.Printf("\nRunning on  Action: %s\n", task.Action)
			cheflinuxoperations.RemoveFile(task.Path)
		default:
			fmt.Println("No action was defined")
		}
	}
	fmt.Println("All Tasks Completed")
}
