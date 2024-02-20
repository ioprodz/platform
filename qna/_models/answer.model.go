package qna_models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AnsweredQuestion struct {
	Question
	AnswerValue int
}

func (aq AnsweredQuestion) IsCorrect() bool {
	var correctAnswer int
	for i, value := range aq.Choices {
		if value.IsCorrect {
			correctAnswer = i
			break
		}
	}
	return aq.AnswerValue == correctAnswer
}

type Answers struct {
	Id   string             `json:"id"`
	Name string             `json:"name"`
	List []AnsweredQuestion `json:"list"`
}

func (ans Answers) MaxScore() int {
	return len(ans.List)
}

func (ans Answers) Score() int {
	var score = 0
	for _, value := range ans.List {
		if value.IsCorrect() {
			score++
		}
	}
	return score
}

func AnswerFromJSON(jsonData []byte) Answers {
	var answer Answers
	if err := json.Unmarshal(jsonData, &answer); err != nil {
		panic("unable to parse QNA json")

	}
	return answer
}

func NewAnswer(qna QNA, answers []int) Answers {

	if len(qna.Questions) != len(answers) {
		panic("lol")
	}

	answerValues := make([]AnsweredQuestion, len(qna.Questions))
	for i, question := range qna.Questions {
		answerValues[i] = AnsweredQuestion{Question: question, AnswerValue: answers[i]}
	}

	return Answers{
		Id:   uuid.NewString(),
		Name: qna.Name,
		List: answerValues,
	}
}
