package members_infra

import (
	"ioprodz/common/db"
	members_models "ioprodz/members/_models"
)

type MemberMongoRepository struct {
	base db.BaseMongoRepository[members_models.Member]
}

func (b *MemberMongoRepository) Create(entity members_models.Member) error {
	return b.base.Create(entity)
}

func (b *MemberMongoRepository) List() ([]members_models.Member, error) {
	return b.base.List()
}

func (b *MemberMongoRepository) Get(id string) (members_models.Member, error) {
	return b.base.Get(id)
}

func (b *MemberMongoRepository) Update(entity members_models.Member) error {
	return b.base.Update(entity)
}

func (b *MemberMongoRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func CreateMongoMemberRepo() *MemberMongoRepository {

	repo := &MemberMongoRepository{base: *db.CreateMongoRepo[members_models.Member]("members")}
	//repo.seed()
	return repo
}

// func (repo *MemberMongoRepository) seed() {

// 	repo.Create(members_models.MemberFromJSON([]byte(`{
// 		"id":"member-id",
// 		"email":"osminosm@gmail.com",
// 		"name":"Osmane Kalache",
// 		"avatarUrl":"https://avatars.githubusercontent.com/u/7093627?v=4",
// 		"accounts":[]
// 	}`)))
// }
