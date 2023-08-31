package helper

import (
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

func CmdRunner(conn *ssh.Client, cmd string) string {
	// Start a session for each command
	var outputBuilder strings.Builder
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Create a multi-writer to capture output for this task
	session.Stdout = io.MultiWriter(log.Writer(), &outputBuilder)
	session.Stderr = io.MultiWriter(log.Writer(), &outputBuilder)

	fmt.Println("Running command:", cmd)
	err = session.Run(cmd)

	// fmt.Println("here", outputBuilder.String())
	if err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}
	return outputBuilder.String()
}
