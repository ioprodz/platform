package qna_infra

import qna_models "ioprodz/qna/_models"

type RepositoryError struct {
	Message string
}

// Error returns the error message
func (e *RepositoryError) Error() string {
	return e.Message
}

type answerRepository struct {
	list []qna_models.Answers
}

func (repo *answerRepository) Create(qna qna_models.Answers) {
	repo.list = append(repo.list, qna)
}

func (repo *answerRepository) List() []qna_models.Answers {
	return repo.list
}

func (repo *answerRepository) Get(id string) (qna_models.Answers, error) {
	for _, obj := range repo.list {
		if obj.Id == id {
			return obj, nil
		}
	}
	return qna_models.Answers{}, &RepositoryError{Message: "Element not found by id: " + id}
}

func CreateAnswerRepo() *answerRepository {
	repo := &answerRepository{list: make([]qna_models.Answers, 0)}
	repo.seed()
	return repo
}

func (repo *answerRepository) seed() {
	repo.list = []qna_models.Answers{qna_models.AnswerFromJSON([]byte(`{"id":"f3c5fe9c-d767-4e30-85eb-43b7c0504585","name":"Object Oriented Programming","list":[{"value":"What is the primary purpose of encapsulation in OOP?","choices":[{"value":"To improve the performance of applications","isCorrect":false},{"value":"To hide the internal state of an object from the outside","isCorrect":true},{"value":"To make code run faster","isCorrect":false},{"value":"To increase the size of the codebase","isCorrect":false}],"AnswerValue":1},{"value":"Which of the following is an example of polymorphism?","choices":[{"value":"A function that adds two numbers","isCorrect":false},{"value":"Using a single function to sort different types of data structures","isCorrect":true},{"value":"Creating multiple instances of a class","isCorrect":false},{"value":"Declaring variables","isCorrect":false}],"AnswerValue":1},{"value":"What does the 'inheritance' concept in OOP allow for?","choices":[{"value":"A class to pass on its properties and methods to another class","isCorrect":true},{"value":"A function to inherit properties from another function","isCorrect":false},{"value":"A method to be executed in parallel","isCorrect":false},{"value":"A class to be duplicated","isCorrect":false}],"AnswerValue":0},{"value":"What does the term 'abstraction' refer to in OOP?","choices":[{"value":"Removing all the functionalities of an object","isCorrect":false},{"value":"Hiding complex implementation details and showing only the necessary features of an object","isCorrect":true},{"value":"The process of creating abstract classes","isCorrect":false},{"value":"The division of a program into smaller programs","isCorrect":false}],"AnswerValue":1},{"value":"Which keyword is used to define an interface in Java?","choices":[{"value":"class","isCorrect":false},{"value":"interface","isCorrect":true},{"value":"implements","isCorrect":false},{"value":"extends","isCorrect":false}],"AnswerValue":3}]}`))}
}
