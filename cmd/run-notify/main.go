package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fako1024/rc-notify"
)

func main() {

	var (
		req            rc.Request
		uri, message   string
		skipSuccessful bool
		printToConsole bool
		returnValue    int
	)

	flag.StringVar(&uri, "uri", "", "RocketChat URI for transmission")
	flag.StringVar(&req.Channel, "chan", "", "Channel to emit to")
	flag.StringVar(&req.User, "user", "", "User to emit as")
	flag.BoolVar(&skipSuccessful, "skip-successful", false, "Skip notification if command was successful")
	flag.BoolVar(&printToConsole, "print-console", false, "Emit messages to console / shell as well")
	flag.Parse()

	// Execute the command based on the remaining command-line arguments
	output, err := runShellCmd(flag.Args())
	if output != "" {
		message += fmt.Sprintln(output)
	}

	// If an error occurred, prepare the notification
	if err != nil {
		message += fmt.Sprintln(err)
		req.Message = fmt.Sprintf("Command `%s` failed:\n```%s```", strings.Join(flag.Args(), " "), message)
		req.Emoji = rc.EmojiWarning
		returnValue = 1
	} else {

		// If successful execution should be notified, prepare the notification
		if !skipSuccessful {
			req.Message = fmt.Sprintf("Command `%s` successful:\n```%s```", strings.Join(flag.Args(), " "), message)
			req.Emoji = rc.EmojiInfo
		}
	}

	// If emission to console was requested, output the result
	if printToConsole {
		fmt.Println(message)
	}

	// If message is present, emit it
	if req.Message != "" {
		if err := emitNotification(uri, req); err != nil {
			fmt.Println(err)
			returnValue = 2
		}
	}

	os.Exit(returnValue)
}

func emitNotification(uri string, req rc.Request) error {

	// Validate the request
	if err := req.Validate(); err != nil {
		return fmt.Errorf("Invalid request: %s", err)
	}

	// Execute the request
	if err := rc.Send(uri, req); err != nil {
		return fmt.Errorf("Failed to send message: %s", err)
	}

	return nil
}

func runShellCmd(args []string) (string, error) {

	if len(args) == 0 {
		return "", nil
	}

	var outBuf bytes.Buffer

	// // Parse command line into command + arguments
	// fields, err := shlex.Split(command)
	// if err != nil || len(fields) == 0 {
	// 	return "", fmt.Errorf("failed to parse command (%s): %s", command, err)
	// }

	// Execute command
	err := generateCommand(args, &outBuf).Run()

	return outBuf.String(), err
}

func generateCommand(fields []string, outBuf io.Writer) (cmd *exec.Cmd) {

	// Check if any arguments were provided
	/* #nosec G204 */
	if len(fields) == 1 {
		cmd = exec.Command(fields[0])
	} else {
		cmd = exec.Command(fields[0], fields[1:]...)
	}

	// Attach STDOUT + STDERR to output buffer
	cmd.Stdout = outBuf
	cmd.Stderr = outBuf

	return
}
