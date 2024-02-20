package qna_models

import "ioprodz/common/policies"

type QNARepository interface {
	policies.Repository[QNA]
}

type AnswersRepository interface {
	policies.Repository[Answers]
}
