package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fako1024/rc-notify"
)

func main() {

	var (
		req  rc.Request
		uri  string
		code bool
	)

	flag.StringVar(&uri, "uri", "", "RocketChat URI for transmission")
	flag.StringVar(&req.Channel, "chan", "", "Channel to emit to")
	flag.StringVar(&req.User, "user", "", "User to emit as")
	flag.StringVar(&req.Message, "msg", "", "Message to send")
	flag.StringVar(&req.Emoji, "emoji", "", "Emoji for the message")
	flag.BoolVar(&code, "code", false, "Emit message as code")
	flag.Parse()

	// Validate the request
	if err := req.Sanitize(); err != nil {
		log.Fatalf("Invalid request: %s", err)
	}

	// Wrap in code markup if requested
	if code {
		req.Message = fmt.Sprintf("```%s```", req.Message)
	}

	// Execute the request
	if err := rc.Send(uri, req); err != nil {
		log.Fatalf("Failed to send message: %s", err)
	}
}
