package cheflinuxoperations

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func stringToFileMode(s string) (os.FileMode, error) {
	mode, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(mode), nil
}

func RemoteFile(localFilePath, remoteFilePath, desiredOwner, desiredGroup string, permissions string) bool {
	perm, err := stringToFileMode(permissions)
	if err != nil {
		log.Fatalf("Failed: %v\n", err)
	}
	desiredPermissions := os.FileMode(perm)
	localContents, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		log.Fatalf("Failed to read local file: %v\n", err)
	}
	out := false
	// Check if the remote file exists
	_, err = os.Stat(remoteFilePath)
	if os.IsNotExist(err) {
		// Create the remote file with the contents of the local file
		err := ioutil.WriteFile(remoteFilePath, localContents, desiredPermissions)
		if err != nil {
			log.Fatalf("Failed to create remote file: %v\n", err)
		} else {
			fmt.Printf("Remote file %s created with desired contents.\n", remoteFilePath)
		}
	} else if err == nil {
		// Read the current contents of the remote file
		remoteContents, err := ioutil.ReadFile(remoteFilePath)
		if err != nil {
			log.Fatalf("Failed to read remote file: %v\n", err)
		}
		if string(remoteContents) != string(localContents) {
			// Update the contents of the remote file to match the local file
			err := ioutil.WriteFile(remoteFilePath, localContents, desiredPermissions)
			if err != nil {
				log.Fatalf("Failed to update remote file contents: %v\n", err)
			} else {
				fmt.Printf("Remote file %s updated with desired contents.\n", remoteFilePath)
				out = true
			}
		} else {
			fmt.Printf("Remote file %s already contains the desired contents.\n", remoteFilePath)
		}

		// Check and update file permissions, owner, and group
		remoteFileStat, err := os.Stat(remoteFilePath)
		if err != nil {
			log.Fatalf("Failed to get remote file info: %v\n", err)
		}

		// Compare permissions, owner, and group
		if remoteFileStat.Mode() != desiredPermissions || getOwnerName(remoteFileStat) != desiredOwner || getGroupName(remoteFileStat) != desiredGroup {
			err := os.Chmod(remoteFilePath, desiredPermissions)
			if err != nil {
				log.Fatalf("Failed to update remote file permissions: %v\n", err)
			} else {
				fmt.Printf("Remote file %s permissions updated.\n", remoteFilePath)
				out = true
			}

			err = os.Chown(remoteFilePath, getOwnerID(desiredOwner), getGroupID(desiredGroup))
			if err != nil {
				log.Fatalf("Failed to update remote file owner/group: %v\n", err)
			} else {
				fmt.Printf("Remote file %s owner/group updated.\n", remoteFilePath)
				out = true
			}
		} else {
			fmt.Printf("Remote file %s permissions, owner, and group are already as desired.\n", remoteFilePath)
		}
	} else {
		log.Fatalf("Error checking remote file status: %v\n", err)
	}
	return out
}

func getOwnerName(info os.FileInfo) string {
	uid := info.Sys().(*syscall.Stat_t).Uid
	user, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return ""
	}
	return user.Username
}

func getGroupName(info os.FileInfo) string {
	gid := info.Sys().(*syscall.Stat_t).Gid
	group, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return ""
	}
	return group.Name
}

func getOwnerID(username string) int {
	user, err := user.Lookup(username)
	if err != nil {
		return -1
	}
	uid, _ := strconv.Atoi(user.Uid)
	return uid
}

func getGroupID(groupname string) int {
	group, err := user.LookupGroup(groupname)
	if err != nil {
		return -1
	}
	gid, _ := strconv.Atoi(group.Gid)
	return gid
}

func RemoveFile(filePath string) {

	// Check if the remote file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("Remote file %s does not exist.\n", filePath)
	} else if err == nil {
		// Remove the remote file
		err := os.Remove(filePath)
		if err != nil {
			log.Fatalf("Failed to remove remote file: %v\n", err)
		} else {
			fmt.Printf("Remote file %s removed successfully.\n", filePath)
		}
	} else {
		log.Fatalf("Error checking remote file status: %v\n", err)
	}
}
