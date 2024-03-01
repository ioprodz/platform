package cv_infra

import (
	"ioprodz/common/db"
	cv_models "ioprodz/cv/_models"
)

type CVMemoryRepository struct {
	base db.BaseMemoryRepository[cv_models.CV]
}

func (b *CVMemoryRepository) Create(entity cv_models.CV) error {
	return b.base.Create(entity)
}

func (b *CVMemoryRepository) List() ([]cv_models.CV, error) {
	return b.base.List()
}

func (b *CVMemoryRepository) Get(id string) (cv_models.CV, error) {
	return b.base.Get(id)
}

func (b *CVMemoryRepository) Update(entity cv_models.CV) error {
	return b.base.Update(entity)
}

func (b *CVMemoryRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func CreateMemoryCVRepo() *CVMemoryRepository {

	repo := &CVMemoryRepository{base: *db.CreateMemoryRepo[cv_models.CV]()}
	repo.seed()
	return repo
}

func (r *CVMemoryRepository) seed() {

	r.Create(cv_models.CVFromJSON([]byte(`{
		"id":"cv-id",
		"userId":"user-id",
		"title":"full stack procrastinator",
		"abstract":"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
		"avatarUrl":"https://i.ytimg.com/vi/JvqJV3bSyaY/hqdefault.jpg?sqp=-oaymwEcCPYBEIoBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==&rs=AOn4CLDp00pLkPy7Qux239R8Ne5qg-uOpA",
		"personal":{
			"fullName":"Djilali Tiarti",
			"email":"djelloul@google.com",
			"phone":"0772923787",
			"address":"123 rue de la paix"
		},
		"education":[
			{
				"period": { "start" : "sep 2023", "end": "sep 2024"},
				"school": "school of rock",
				"degree": "phd"
			}
		],
		"experience":[
			{
				"period":{ "start" : "sep 2023", "end": "sep 2024"},
				"company":"SNTA",
				"title":"general manager",
				"details":"- wake up \n-go to office\n -do nothing \n - go back home \n -repeat"
			}
		]
	}`)))
}
