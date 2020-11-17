# A Simple tool to emit messages via RocketChat Webhook REST API

[![Github Release](https://img.shields.io/github/release/fako1024/rc-notify.svg)](https://github.com/fako1024/rc-notify/releases)
[![GoDoc](https://godoc.org/github.com/fako1024/rc-notify?status.svg)](https://godoc.org/github.com/fako1024/rc-notify/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fako1024/rc-notify)](https://goreportcard.com/report/github.com/fako1024/rc-notify)
[![Build/Test Status](https://github.com/fako1024/rc-notify/workflows/Go/badge.svg)](https://github.com/fako1024/rc-notify/actions?query=workflow%3AGo)

This package provides an interface to RocketChat's Webhook / notification REST API (provided an endpoint is made available) and a command line tool for simple message sending / delivery.

## Installation
```bash
go get -u github.com/fako1024/rc-notify
```

## Examples
#### Send a simple message
```go
req := rc.Request{
	Channel: "@me",
	User:    "Sending User",
	Message: "My message",
	Emoji:   rc.EmojiAlert,
}

// Validate the request
if err := req.Validate(); err != nil {
	log.Fatalf("Invalid request: %s", err)
}

// Execute the request (set webhook URI accordingly)
uri := "https://your.rc.instance.com/hooks/.../..."
if err := rc.Send(uri, req); err != nil {
	log.Fatalf("Failed to send message: %s", err)
}
```
