package members_infra

import (
	"ioprodz/common/db"
	members_models "ioprodz/members/_models"
)

type UserProjectMongoRepository struct {
	base db.BaseMongoRepository[members_models.UserProject]
}

func (b *UserProjectMongoRepository) Create(entity members_models.UserProject) error {
	return b.base.Create(entity)
}

func (b *UserProjectMongoRepository) List() ([]members_models.UserProject, error) {
	return b.base.List()
}

func (b *UserProjectMongoRepository) Get(id string) (members_models.UserProject, error) {
	return b.base.Get(id)
}

func (b *UserProjectMongoRepository) Update(entity members_models.UserProject) error {
	return b.base.Update(entity)
}

func (b *UserProjectMongoRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (b *UserProjectMongoRepository) GetByUserId(userId string) ([]members_models.UserProject, error) {
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

func CreateMongoUserProjectRepo() *UserProjectMongoRepository {

	repo := &UserProjectMongoRepository{base: *db.CreateMongoRepo[members_models.UserProject]("user_projects")}
	//repo.seed()
	return repo
}

// func (repo *UserProjectMongoRepository) seed() {

// 	repo.Create(members_models.UserProjectFromJSON([]byte(`{
// 		"id":"member-id",
// 		"email":"osminosm@gmail.com",
// 		"name":"Osmane Kalache",
// 		"avatarUrl":"https://avatars.githubusercontent.com/u/7093627?v=4",
// 		"accounts":[]
// 	}`)))
// }
