package qna_models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type QuestionChoice struct {
	Value     string `json:"value"`
	IsCorrect bool   `json:"isCorrect"`
}

type Question struct {
	Value   string           `json:"value"`
	Choices []QuestionChoice `json:"choices"`
}

type QNA struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
}

func (qna *QNA) SetQuestions(questions []Question) {
	qna.Questions = questions
}

func QNAFromJSON(jsonData []byte) QNA {
	var qna QNA
	if err := json.Unmarshal(jsonData, &qna); err != nil {
		panic("unable to parse QNA json")

	}
	return qna
}

func NewQNA(name string) *QNA {
	return &QNA{
		Id:        uuid.NewString(),
		Name:      name,
		Questions: make([]Question, 0),
	}
}
