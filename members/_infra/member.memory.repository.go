package members_infra

import (
	"ioprodz/common/db"
	members_models "ioprodz/members/_models"
)

type MemberMemoryRepository struct {
	base db.BaseMemoryRepository[members_models.Member]
}

func (b *MemberMemoryRepository) Create(entity members_models.Member) error {
	return b.base.Create(entity)
}

func (b *MemberMemoryRepository) List() ([]members_models.Member, error) {
	return b.base.List()
}

func (b *MemberMemoryRepository) Get(id string) (members_models.Member, error) {
	return b.base.Get(id)
}

func (b *MemberMemoryRepository) Update(entity members_models.Member) error {
	return b.base.Update(entity)
}

func (b *MemberMemoryRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func CreateMemoryMemberRepo() *MemberMemoryRepository {

	repo := &MemberMemoryRepository{base: *db.CreateMemoryRepo[members_models.Member]()}
	//repo.seed()
	return repo
}

// func (repo *MemberMemoryRepository) seed() {

// 	repo.Create(members_models.MemberFromJSON([]byte(`{
// 		"id":"member-id",
// 		"email":"osminosm@gmail.com",
// 		"name":"Osmane Kalache",
// 		"bio":"I love computers",
// 		"avatarUrl":"https://avatars.githubusercontent.com/u/7093627?v=4",
// 		"accounts":[],
// 		"links":[
// 			{ "name":"github", "url":"https://github.com/osminosm" },
// 			{ "name":"linkedin", "url":"https://linkedin.com/osminosm" }
// 		]
// 	}`)))
// }
