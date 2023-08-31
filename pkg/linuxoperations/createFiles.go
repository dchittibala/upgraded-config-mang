package linuxoperations

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/reddydinesh427/upgraded-config-mang/pkg/helper"
	"golang.org/x/crypto/ssh"
)

func CreateDirectory(conn *ssh.Client, directoryPath string, perm, owner, group string) {
	checkServiceStatusCmd := fmt.Sprintf("[ -d %s ] && echo 'exists' || echo 'not exists'", directoryPath)
	serviceStatus := helper.CmdRunner(conn, checkServiceStatusCmd)
	if owner == "" {
		owner = helper.CmdRunner(conn, "id -u")
	}
	if group == "" {
		group = helper.CmdRunner(conn, "id -g")
	}
	if strings.Contains(serviceStatus, "not exists") {
		fmt.Println("create directory:", directoryPath)
		startServiceCmd := fmt.Sprintf("sudo mkdir %s;sudo chmod %s %s;sudo chown %s:%s %s", directoryPath, perm, directoryPath, owner, group, directoryPath)
		helper.CmdRunner(conn, startServiceCmd)
	} else {
		fmt.Println("directory exists,validating perms")
		checkDirDetailsCmd := fmt.Sprintf("[[ $(stat -c '%%a %%U %%G' %s) = \"%s %s %s\" ]] && echo 'correct' || echo 'incorrect'", directoryPath, perm, owner, group)
		if helper.CmdRunner(conn, checkDirDetailsCmd) == "incorrect" {
			startServiceCmd := fmt.Sprintf("sudo chmod %s %s && sudo chown %s:%s %s", perm, directoryPath, owner, group, directoryPath)
			helper.CmdRunner(conn, startServiceCmd)
		} else {
			fmt.Println("No change in directory permisions")
		}
	}
}

// func fileModeToString(fileMode os.FileMode) string {
// 	permissions := fileMode.Perm()
// 	permString := fmt.Sprintf("%03o", permissions)
// 	return permString
// }

func RemoteFileCopy(conn *ssh.Client, localFilePath, remoteFilePath, owner, group, perm string) bool {
	if perm == "" || localFilePath == "" || remoteFilePath == "" || owner == "" || group == "" {
		panic("One or more Key values are empty, Please check the name of the key for localPath, remotePath, owner, group, perm for action: remoteFile in tasks.json")
	}
	tempCopy := CopyFileContents(conn, localFilePath)
	// remoteFileCheckCmd := fmt.Sprintf("[ -f %s ] && cmp -s %s - < %s || echo 'different'", remoteFilePath, remoteFilePath, remoteFilePath)
	remoteFileCheckCmd := fmt.Sprintf("[ -f %s ] && cmp -s %s %s || echo 'different'", remoteFilePath, remoteFilePath, tempCopy)

	serviceStatus := helper.CmdRunner(conn, remoteFileCheckCmd)
	// if serviceStatus == "different" {
	if strings.Contains(serviceStatus, "different") {
		fmt.Printf("remote file and local file are different updating the remote file %s with %s\n", remoteFilePath, localFilePath)
		moveFile := fmt.Sprintf("mv %s %s", tempCopy, remoteFilePath)
		helper.CmdRunner(conn, moveFile)
	}
	// delete the temperory file created
	deleteFile := fmt.Sprintf("rm -rf %s", tempCopy)
	helper.CmdRunner(conn, deleteFile)

	checkfileDetailsCmd := fmt.Sprintf("[[ $(stat -c '%%a %%U %%G' %s) = \"%s %s %s\" ]] && echo 'correct' || echo 'incorrect'", remoteFilePath, perm, owner, group)
	if strings.Contains(helper.CmdRunner(conn, checkfileDetailsCmd), "incorrect") {
		fmt.Printf("Incorrect File Permissions,updating the file perms for %s\n", remoteFilePath)
		startServiceCmd := fmt.Sprintf("sudo chmod %s %s && sudo chown %s:%s %s", perm, remoteFilePath, owner, group, remoteFilePath)
		helper.CmdRunner(conn, startServiceCmd)
	}
	return strings.Contains(serviceStatus, "different")
}

func CopyFileContents(conn *ssh.Client, localFilePath string) string {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()
	// remoteFile, err := session.StdinPipe()
	// Set up pipes for input/output
	var outputBuilder strings.Builder
	session.Stdout = io.MultiWriter(log.Writer(), &outputBuilder)
	session.Stderr = io.MultiWriter(log.Writer(), &outputBuilder)

	// Open the local file
	localFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatalf("Failed to open local file: %v", err)
	}
	defer localFile.Close()
	var localFileContent bytes.Buffer
	_, err = io.Copy(&localFileContent, localFile)
	if err != nil {
		log.Fatalf("Failed to read local file contents: %v", err)
	}
	remoteFile, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to open remote file for writing: %v", err)
	}
	defer remoteFile.Close()
	// currentTime := time.Now()
	// formattedTime := currentTime.Format("2006-01-02 15:04:05")
	tempFile := "/tmp/temp" + strconv.FormatInt(time.Now().Unix(), 10)
	err = session.Start(fmt.Sprintf("cat > %s", tempFile))
	if err != nil {
		log.Fatalf("Failed to start remote command: %v", err)
	}

	// Copy the contents of the local file to the remote file
	_, err = io.Copy(remoteFile, &localFileContent)
	if err != nil {
		log.Fatalf("Failed to copy contents to remote file: %v", err)
	}
	// Close the remote file and wait for the remote command to finish
	err = remoteFile.Close()
	if err != nil {
		log.Fatalf("Failed to close remote file: %v", err)
	}
	err = session.Wait()
	if err != nil {
		log.Fatalf("Remote command failed: %v", err)
	}
	return tempFile
}

func RemoveFile(conn *ssh.Client, filePath string) {
	if filePath == "" {
		panic("filePath is empty,  action: removeFile in tasks.json")
	}
	remoteFileCheckCmd := fmt.Sprintf("[ -f %s ] && echo 'exists' || echo 'not_exists'", filePath)
	serviceStatus := helper.CmdRunner(conn, remoteFileCheckCmd)
	if strings.Contains(serviceStatus, "not_exists") {
		fmt.Println("File doesnt Exists:", filePath)
	} else {
		fmt.Println("File exists,removing file:", filePath)
		checkDirDetailsCmd := fmt.Sprintf("rm -rf %s", filePath)
		helper.CmdRunner(conn, checkDirDetailsCmd)
	}
}
