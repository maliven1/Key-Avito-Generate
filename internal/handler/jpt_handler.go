package handler

import (
	"avito/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const key = ""

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func JptHandler(keys string, titleKey string, i int) {
	const op = "internal.server.JptHandler"
	log := logger.Log.With(
		slog.String("op", op),
	)
	log.Info("Starting JPT Handler")

	url := "https://api.intelligence.io.solutions/api/v1/chat/completions"

	request := Request{
		Model: "deepseek-ai/DeepSeek-R1",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant",
			},
			{
				Role:    "user",
				Content: "оставь только 60% ключей, которые подходят по тематике данных ключей" + titleKey + ". Выбирай из этих ключей" + keys + ".В ответ отправь только выбранные ключи.",
			},
		},
	}

	payload := strings.NewReader(fmt.Sprintf(`{"model":"%s","messages":[{"role":"%s","content":"%s"},{"role":"%s","content":"%s"}]}`,
		request.Model,
		request.Messages[0].Role,
		request.Messages[0].Content,
		request.Messages[1].Role,
		request.Messages[1].Content,
	))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer func() {
		err := res.Body.Close()
		if err != nil {
			logger.Log.Error("Error closing response body", err)
		}
	}()
	body, _ := io.ReadAll(res.Body)

	var response Response
	json.Unmarshal(body, &response)

	if len(response.Choices) > 0 {
		content := response.Choices[0].Message.Content
		if idx := strings.Index(content, "</think>"); idx != -1 {
			content = content[idx+8:]
		}
		fmt.Println(i, titleKey+"\n"+content)
	}

}
