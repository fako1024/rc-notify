package rc

import (
	"fmt"

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

	if err := r.Validate(); err != nil {
		return fmt.Errorf("Error validating RocketChat request: %s", err)
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

// Validate checks the required fields of a request
func (r Request) Validate() error {
	if r.Channel == "" {
		return fmt.Errorf("Channel parameter missing")
	}
	if r.Message == "" {
		return fmt.Errorf("Message parameter missing")
	}
	return nil
}