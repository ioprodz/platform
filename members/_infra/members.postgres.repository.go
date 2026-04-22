package members_infra

import (
	"ioprodz/common/db"
	members_models "ioprodz/members/_models"
)

type MemberPostgresRepository struct {
	base db.BasePostgresRepository[members_models.Member]
}

func (b *MemberPostgresRepository) Create(entity members_models.Member) error {
	return b.base.Create(entity)
}

func (b *MemberPostgresRepository) List() ([]members_models.Member, error) {
	return b.base.List()
}

func (b *MemberPostgresRepository) Get(id string) (members_models.Member, error) {
	return b.base.Get(id)
}

func (b *MemberPostgresRepository) Update(entity members_models.Member) error {
	return b.base.Update(entity)
}

func (b *MemberPostgresRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func CreatePostgresMemberRepo() *MemberPostgresRepository {

	repo := &MemberPostgresRepository{base: *db.CreatePostgresRepo[members_models.Member]("members")}
	return repo
}
