package members_infra

import (
	"ioprodz/common/db"
	members_models "ioprodz/members/_models"
)

type UserProjectPostgresRepository struct {
	base db.BasePostgresRepository[members_models.UserProject]
}

func (b *UserProjectPostgresRepository) Create(entity members_models.UserProject) error {
	return b.base.Create(entity)
}

func (b *UserProjectPostgresRepository) List() ([]members_models.UserProject, error) {
	return b.base.List()
}

func (b *UserProjectPostgresRepository) Get(id string) (members_models.UserProject, error) {
	return b.base.Get(id)
}

func (b *UserProjectPostgresRepository) Update(entity members_models.UserProject) error {
	return b.base.Update(entity)
}

func (b *UserProjectPostgresRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *UserProjectPostgresRepository) GetByUserId(userId string) ([]members_models.UserProject, error) {
	list, err := b.List()
	result := []members_models.UserProject{}
	if err != nil {
		return result, err
	}
	for _, porject := range list {
		if porject.UserId == userId {
			result = append(result, porject)
		}
	}

	return result, nil
}

func CreatePostgresUserProjectRepo() *UserProjectPostgresRepository {

	repo := &UserProjectPostgresRepository{base: *db.CreatePostgresRepo[members_models.UserProject]("user_projects")}
	return repo
}
