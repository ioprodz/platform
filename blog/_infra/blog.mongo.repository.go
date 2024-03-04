package blog_infra

import (
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/db"
)

type BlogMongoRepository struct {
	base db.BaseMongoRepository[blog_models.Blog]
}

func (b *BlogMongoRepository) Create(entity blog_models.Blog) error {
	return b.base.Create(entity)
}

func (b *BlogMongoRepository) List() ([]blog_models.Blog, error) {
	return b.base.List()
}

func (b *BlogMongoRepository) Get(id string) (blog_models.Blog, error) {
	return b.base.Get(id)
}

func (b *BlogMongoRepository) Update(entity blog_models.Blog) error {
	return b.base.Update(entity)
}

func (b *BlogMongoRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (repo *BlogMongoRepository) ListPublished() ([]blog_models.Blog, error) {
	list := []blog_models.Blog{}
	allPosts, _ := repo.List()
	for _, post := range allPosts {
		if post.IsPublished() {
			list = append(list, post)
		}
	}
	return list, nil
}

func CreateMongoBlogRepo() *BlogMongoRepository {

	repo := &BlogMongoRepository{base: *db.CreateMongoRepo[blog_models.Blog]("blogposts")}
	return repo
}
