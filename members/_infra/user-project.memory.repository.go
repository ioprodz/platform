package members_infra

import (
	"ioprodz/common/db"
	members_models "ioprodz/members/_models"
)

type UserProjectMemoryRepository struct {
	base db.BaseMemoryRepository[members_models.UserProject]
}

func (b *UserProjectMemoryRepository) Create(entity members_models.UserProject) error {
	return b.base.Create(entity)
}

func (b *UserProjectMemoryRepository) List() ([]members_models.UserProject, error) {
	return b.base.List()
}

func (b *UserProjectMemoryRepository) Get(id string) (members_models.UserProject, error) {
	return b.base.Get(id)
}

func (b *UserProjectMemoryRepository) Update(entity members_models.UserProject) error {
	return b.base.Update(entity)
}

func (b *UserProjectMemoryRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *UserProjectMemoryRepository) GetByUserId(userId string) ([]members_models.UserProject, error) {
	list, err := b.List()
	result := []members_models.UserProject{}
	if err != nil {
		return result, err
	}
	for _, project := range list {
		if project.UserId == userId {
			result = append(result, project)
		}
	}

	return result, nil
}

func CreateMemoryUserProjectRepo() *UserProjectMemoryRepository {

	repo := &UserProjectMemoryRepository{base: *db.CreateMemoryRepo[members_models.UserProject]()}
	repo.seed()
	return repo
}

func (repo *UserProjectMemoryRepository) seed() {

	repo.Create(members_models.UserProjectFromJSON([]byte(`{
		"id":"project-id",
		"title":"cool project",
		"userId":"user-id",
		"description":"its the best project ever"
	}`)))
}
