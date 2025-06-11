package openaiClient

import (
	"encoding/json"
	"fmt"
	"ioprodz/common/clients/httpClients"
)

type JsonPromptChoiceMessage struct {
	Content string `json:"content"`
}
type JsonPromptChoice struct {
	Message JsonPromptChoiceMessage `json:"message"`
}

type JsonPromptResponse struct {
	Choices []JsonPromptChoice `json:"choices"`
}

func JsonPrompt(input string, jsonReponseFormat string) (string, error) {
	response, err := httpClients.Post("https://api.openai.com/v1/chat/completions", map[string]interface{}{
		"model":           "gpt-4o",
		"response_format": map[string]string{"type": "json_object"},
		"messages":        []map[string]string{{"role": "user", "content": input + " in json format: " + jsonReponseFormat}},
		"temperature":     0.7,
	})
	if err != nil {
		fmt.Println("Error calling the open ai api err:", err)
		return "", err
	}

	var jsonReponse JsonPromptResponse
	if err := json.Unmarshal(response, &jsonReponse); err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return "", err
	}

	return jsonReponse.Choices[0].Message.Content, nil
}
