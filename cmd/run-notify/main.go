package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fako1024/rc-notify"
	"github.com/sevlyar/go-daemon"
)

const (
	errCommand   = 1
	errEmission  = 2
	errDaemonize = 3
)

func main() {

	var (
		ctx             *daemon.Context
		req             rc.Request
		uri, message    string
		background      bool
		skipSuccessful  bool
		printToConsole  bool
		returnValue     int
		messageMaxLines int
	)

	flag.StringVar(&uri, "uri", "", "RocketChat URI for transmission")
	flag.StringVar(&req.Channel, "chan", "", "Channel to emit to")
	flag.StringVar(&req.User, "user", "", "User to emit as")
	flag.BoolVar(&background, "background", false, "Run command in the background")
	flag.BoolVar(&skipSuccessful, "skip-successful", false, "Skip notification if command was successful")
	flag.BoolVar(&printToConsole, "print-console", false, "Emit messages to console / shell as well")
	flag.IntVar(&messageMaxLines, "max-lines", -1, "Maximum number of command output lines to emit (default: unlimited)")
	flag.Parse()

	// If requested, daemonize / run in background
	if background {
		ctx = new(daemon.Context)

		d, err := ctx.Reborn()
		handleErr(err, 3)
		if d != nil {
			return
		}
		defer handleErr(ctx.Release(), errDaemonize)
	}

	// Execute the command based on the remaining command-line arguments
	output, err := runShellCmd(flag.Args())
	if output != "" {
		message += fmt.Sprintln(output)
	}

	// If an error occurred, prepare the notification
	if err != nil {
		message += fmt.Sprintln(err)
		req.Message = fmt.Sprintf("Command `%s` failed:\n```%s```", strings.Join(flag.Args(), " "), limitMessage(message, messageMaxLines))
		req.Emoji = rc.EmojiWarning
		returnValue = errCommand
	} else {

		// If successful execution should be notified, prepare the notification
		if !skipSuccessful {
			req.Message = fmt.Sprintf("Command `%s` successful:\n```%s```", strings.Join(flag.Args(), " "), limitMessage(message, messageMaxLines))
			req.Emoji = rc.EmojiInfo
		}
	}

	// If emission to console was requested, output the result
	if printToConsole {
		fmt.Println(message)
	}

	// If message is present, emit it
	if req.Message != "" {
		handleErr(emitNotification(uri, req), errEmission)
	}

	os.Exit(returnValue)
}

func emitNotification(uri string, req rc.Request) error {

	// Validate the request
	if err := req.Sanitize(); err != nil {
		return fmt.Errorf("invalid request: %s", err)
	}

	// Execute the request
	if err := rc.Send(uri, req); err != nil {
		return fmt.Errorf("failed to send message: %s", err)
	}

	return nil
}

func runShellCmd(args []string) (string, error) {

	if len(args) == 0 {
		return "", nil
	}

	// Execute command
	var outBuf bytes.Buffer
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

func limitMessage(input string, maxLen int) (output string) {

	if maxLen <= 0 {
		return input
	}

	scanner := bufio.NewScanner(strings.NewReader(input))
	linesRead := 0
	for scanner.Scan() {
		output += fmt.Sprintln(scanner.Text())
		linesRead++

		if linesRead >= maxLen {
			return
		}
	}
	return
}

func handleErr(err error, returnValue int) {
	if err == nil {
		return
	}

	fmt.Println(err)
	flag.Usage()
	os.Exit(returnValue)
}
