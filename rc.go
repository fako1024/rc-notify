package rc

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/fako1024/httpc"
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
	Channel     string       `json:"channel"`               // The channel to send to
	User        string       `json:"username"`              // The username to display next to the message
	Alias       string       `json:"alias,omitempty"`       // An alias to add to the message
	Emoji       string       `json:"emoji"`                 // The icon to use for the message
	Message     string       `json:"text"`                  // The message to send
	Attachments []Attachment `json:"attachments,omitempty"` // Potential message attachments
}

// Attachment denotes a message attachment, like subfields or images and links
type Attachment struct {
	Color      string  `json:"color,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Title      string  `json:"title"`
	TitleLink  string  `json:"title_link"`
	Text       string  `json:"text"`
	ImageURL   string  `json:"image_url,omitempty"`
	ThumbURL   string  `json:"thumb_url,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

// Field denotes a (sub-)field to be displayed in the RC message / attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Send sends the request to the defined endpoint / RC instance
func Send(uri string, r Request) error {

	// Validate the request
	if err := r.Sanitize(); err != nil {
		return fmt.Errorf("error validating RocketChat request: %w", err)
	}

	// Prepare and run the request
	return httpc.New("POST", uri).
		Transport(http.DefaultTransport).
		RetryBackOff(httpc.Intervals{
			time.Second,
			5 * time.Second,
		}).
		EncodeJSON(r).
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

type FileUploadRequest struct {
	Data        []byte
	RoomID      string
	Message     string
	Description string
}

type APIAuth struct {
	UserID string
	Token  string
}

func UploadFile(endpoint string, auth APIAuth, req FileUploadRequest) error {

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	if err := writer.SetBoundary("45b03ac0dfd03bafd94e05b0547ed86c5bfb46a201451552a604d6a3aac1"); err != nil {
		return fmt.Errorf("failed to set multipart boundary: %w", err)
	}
	imageFile, err := createImageFormFile(writer, "image.jpg")
	if err != nil {
		return err
	}

	if n, err := imageFile.Write(req.Data); err != nil || n != len(req.Data) {
		return fmt.Errorf("failed to write image data to buffer (err: %w)", err)
	}

	if req.Message != "" {
		if err := writer.WriteField("msg", req.Message); err != nil {
			return fmt.Errorf("failed to set message: %w", err)
		}
	}
	if req.Description != "" {
		if err := writer.WriteField("description", req.Description); err != nil {
			return fmt.Errorf("failed to set description: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return err
	}

	// Prepare and run the request
	return httpc.New("POST", strings.TrimRight(endpoint, "/")+"/api/v1/rooms.upload/"+req.RoomID).
		Transport(http.DefaultTransport).
		RetryBackOff(httpc.Intervals{
			time.Second,
			5 * time.Second,
		}).
		Headers(httpc.Params{
			"X-User-Id":    auth.UserID,
			"X-Auth-Token": auth.Token,
			"Content-Type": writer.FormDataContentType(),
		}).
		Body(body.Bytes()).
		Run()
}

func createImageFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", "image/jpeg")
	return w.CreatePart(h)
}
