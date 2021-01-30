package rc

import (
	"fmt"
	"strings"

	"github.com/fako1024/httpc"
	jsoniter "github.com/json-iterator/go"
)

const (

	// EmojiInfo denotes an info logo
	EmojiInfo = ":information_source:"

	// EmojiWarning denotes a warning logo
	EmojiWarning = ":warning:"

	// EmojiAlert denotes an alert logo
	EmojiAlert = ":rotating_light:"

	// EmojiBulb denotes a lightbulb logo
	EmojiBulb = ":bulb:"

	// EmojiArrowRight denotes a right arrow logo
	EmojiArrowRight = ":arrow_right:"

	// EmojiArrowLeft denotes a left arrow logo
	EmojiArrowLeft = ":arrow_left:"

	// EmojiLoginTest denotes the login logo for TEST
	EmojiLoginTest = ":login_test:"

	// EmojiLogoutTest denotes the logout logo for TEST
	EmojiLogoutTest = ":logout_test:"

	// EmojiLoginProd denotes the login logo for PROD
	EmojiLoginProd = ":login_prod:"

	// EmojiLogoutProd denotes the logout logo for PROD
	EmojiLogoutProd = ":logout_prod:"
)

// Request denotes an RC JSON request
type Request struct {
	Channel string `json:"channel"`    // The channel to send to
	User    string `json:"username"`   // The username to display next to the message
	Emoji   string `json:"icon_emoji"` // The icon to use for the message
	Message string `json:"text"`       // The message to send
}

// Send sends the request to the defined endpoint / RC instance
func Send(uri string, r Request) error {

	// Validate the request
	if err := r.Sanitize(); err != nil {
		return fmt.Errorf("error validating RocketChat request: %s", err)
	}

	// Marshal the request into a JSON structure
	json, err := jsoniter.Marshal(r)
	if err != nil {
		return err
	}

	// Prepare and run the request
	return httpc.New("POST", uri).
		Headers(httpc.Params{
			"Content-Type": "application/json",
		}).
		Body(json).
		Run()
}

// Sanitize checks and sanitizes the required fields of a request
func (r *Request) Sanitize() error {
	if r.Channel == "" {
		return fmt.Errorf("channel parameter missing")
	}
	if r.Message == "" {
		return fmt.Errorf("message parameter missing")
	}

	// If a channel name was provided without any prefix, assume a standard channel
	if !strings.HasPrefix(r.Channel, "#") && !strings.HasPrefix(r.Channel, "@") {
		r.Channel = "#" + r.Channel
	}

	// Set an informational emoji (default would we a warning), if empty
	if r.Emoji == "" {
		r.Emoji = EmojiInfo
	}

	return nil
}
