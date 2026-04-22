package cv_infra

import (
	"ioprodz/common/db"
	cv_models "ioprodz/cv/_models"
)

type CVPostgresRepository struct {
	base db.BasePostgresRepository[cv_models.CV]
}

func (b *CVPostgresRepository) Create(entity cv_models.CV) error {
	return b.base.Create(entity)
}

func (b *CVPostgresRepository) List() ([]cv_models.CV, error) {
	return b.base.List()
}

func (b *CVPostgresRepository) Get(id string) (cv_models.CV, error) {
	return b.base.Get(id)
}

func (b *CVPostgresRepository) Update(entity cv_models.CV) error {
	return b.base.Update(entity)
}

func (b *CVPostgresRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func CreatePostgresCVRepo() *CVPostgresRepository {

	repo := &CVPostgresRepository{base: *db.CreatePostgresRepo[cv_models.CV]("curriculums")}

	return repo
}
