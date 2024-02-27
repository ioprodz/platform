package qna_models

import "ioprodz/common/policies"

type QNARepository interface {
	policies.BaseRepository[QNA]
}

type AnswersRepository interface {
	policies.BaseRepository[Answers]
}
