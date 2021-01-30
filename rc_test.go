package rc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

var (
	uri = "https://chat.example.org/hooks/UnpNPC3DWR0tMy"

	testMessage = Request{
		Channel: "@test",
		User:    "user",
		Emoji:   EmojiAlert,
		Message: "Hello, world!",
	}
	invalidMessageNoChannel = Request{
		Channel: "",
		User:    "user",
		Emoji:   EmojiAlert,
		Message: "Hello, world!",
	}
	invalidMessageNoMessage = Request{
		Channel: "@test",
		User:    "user",
		Emoji:   EmojiAlert,
		Message: "",
	}
)

func TestValidate(t *testing.T) {
	if err := testMessage.Sanitize(); err != nil {
		t.Fatalf("Failed to validate message: %s", err)
	}
	if err := invalidMessageNoChannel.Sanitize(); err == nil || err.Error() != "channel parameter missing" {
		t.Fatalf("Unexpected success when validating message without channel: %s", err)
	}
	if err := invalidMessageNoMessage.Sanitize(); err == nil || err.Error() != "message parameter missing" {
		t.Fatalf("Unexpected success when validating message without message text: %s", err)
	}
}

func TestSanitizeMissingChannelPrefix(t *testing.T) {

	// Set up a mock matcher
	defer gock.Off()
	gock.New(uri).
		Post(path.Base(uri)).
		MatchType("application/json").
		AddMatcher(gock.MatchFunc(func(arg1 *http.Request, arg2 *gock.Request) (bool, error) {
			data, err := ioutil.ReadAll(arg1.Body)
			if err != nil {
				return false, err
			}

			return bytes.Equal(data, []byte(`{"channel":"#test","username":"user","icon_emoji":":rotating_light:","text":"Hello, world!"}`)), nil
		})).
		Reply(http.StatusOK)

	if err := Send(uri, Request{
		Channel: "test",
		User:    "user",
		Message: "Hello, world!",
		Emoji:   EmojiAlert,
	}); err != nil {
		t.Fatalf("Failed to send message: %s", err)
	}
}

func TestSanitizeDefaultEmoji(t *testing.T) {

	// Set up a mock matcher
	defer gock.Off()
	gock.New(uri).
		Post(path.Base(uri)).
		MatchType("application/json").
		AddMatcher(gock.MatchFunc(func(arg1 *http.Request, arg2 *gock.Request) (bool, error) {
			data, err := ioutil.ReadAll(arg1.Body)
			if err != nil {
				return false, err
			}

			return bytes.Equal(data, []byte(`{"channel":"@test","username":"user","icon_emoji":":information_source:","text":"Hello, world!"}`)), nil
		})).
		Reply(http.StatusOK)

	if err := Send(uri, Request{
		Channel: "@test",
		User:    "user",
		Message: "Hello, world!",
	}); err != nil {
		t.Fatalf("Failed to send message: %s", err)
	}
}

func TestSendSimpleMessage(t *testing.T) {

	// Set up a mock matcher
	defer gock.Off()
	gock.New(uri).
		Post(path.Base(uri)).
		MatchType("application/json").
		AddMatcher(gock.MatchFunc(func(arg1 *http.Request, arg2 *gock.Request) (bool, error) {
			data, err := ioutil.ReadAll(arg1.Body)
			if err != nil {
				return false, err
			}

			return bytes.Equal(data, []byte(`{"channel":"@test","username":"user","icon_emoji":":rotating_light:","text":"Hello, world!"}`)), nil
		})).
		Reply(http.StatusOK)

	if err := Send(uri, testMessage); err != nil {
		t.Fatalf("Failed to send message: %s", err)
	}
}

func TestSendInvalidMessage(t *testing.T) {
	if err := Send(uri, invalidMessageNoChannel); err == nil || err.Error() != "error validating RocketChat request: channel parameter missing" {
		t.Fatalf("Unexpected success when sending message without channel: %s", err)
	}
	if err := Send(uri, invalidMessageNoMessage); err == nil || err.Error() != "error validating RocketChat request: message parameter missing" {
		t.Fatalf("Unexpected success when sending message without message text: %s", err)
	}
}
