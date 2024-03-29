package rc

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
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
			data, err := io.ReadAll(arg1.Body)
			if err != nil {
				return false, err
			}

			return bytes.Equal(data, []byte(`{"channel":"#test","username":"user","emoji":":rotating_light:","text":"Hello, world!"}`)), nil
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
			data, err := io.ReadAll(arg1.Body)
			if err != nil {
				return false, err
			}

			return bytes.Equal(data, []byte(`{"channel":"@test","username":"user","emoji":":information_source:","text":"Hello, world!"}`)), nil
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
			data, err := io.ReadAll(arg1.Body)
			if err != nil {
				return false, err
			}

			return bytes.Equal(data, []byte(`{"channel":"@test","username":"user","emoji":":rotating_light:","text":"Hello, world!"}`)), nil
		})).
		Reply(http.StatusOK)

	if err := Send(uri, testMessage); err != nil {
		t.Fatalf("Failed to send message: %s", err)
	}
}

func TestUploadFile(t *testing.T) {

	uploadURI := "https://chat.example.org/api/v1/rooms.upload/randomRoomID"

	// Set up a mock matcher
	defer gock.Off()
	gock.New(uploadURI).
		Post("/api/v1/rooms.upload/randomRoomID").
		MatchType("multipart/form-data; boundary=45b03ac0dfd03bafd94e05b0547ed86c5bfb46a201451552a604d6a3aac1").
		MatchHeaders(map[string]string{
			"X-User-Id":    "testID",
			"X-Auth-Token": "testToken",
		}).
		AddMatcher(gock.MatchFunc(func(arg1 *http.Request, arg2 *gock.Request) (bool, error) {
			reader := multipart.NewReader(arg1.Body, "45b03ac0dfd03bafd94e05b0547ed86c5bfb46a201451552a604d6a3aac1")

			part, err := reader.NextPart()
			if err != nil {
				return false, err
			}

			data, err := io.ReadAll(part)
			if err != nil {
				return false, err
			}

			formData, err := reader.ReadForm(1024)
			if err != nil {
				return false, err
			}
			if len(formData.Value["msg"]) != 1 || formData.Value["msg"][0] != "Test Message" {
				return false, fmt.Errorf("missing / invalid message: %v", formData.Value)
			}
			if len(formData.Value["description"]) != 1 || formData.Value["description"][0] != "Test Description" {
				return false, fmt.Errorf("missing / invalid description: %v", formData.Value)
			}

			return bytes.Equal(data, []byte(`This is a simple text file`)), nil
		})).
		Reply(http.StatusOK)

	if err := UploadFile("https://chat.example.org/", APIAuth{
		UserID: "testID",
		Token:  "testToken",
	}, FileUploadRequest{
		Data:        []byte(`This is a simple text file`),
		RoomID:      "randomRoomID",
		Message:     "Test Message",
		Description: "Test Description",
	}); err != nil {
		t.Fatalf("Failed to upload file: %s", err)
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
