package blog_infra

import (
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/db"
)

type BlogPostgresRepository struct {
	base db.BasePostgresRepository[blog_models.Blog]
}

func (b *BlogPostgresRepository) Create(entity blog_models.Blog) error {
	return b.base.Create(entity)
}

func (b *BlogPostgresRepository) List() ([]blog_models.Blog, error) {
	return b.base.List()
}

func (b *BlogPostgresRepository) Get(id string) (blog_models.Blog, error) {
	return b.base.Get(id)
}

func (b *BlogPostgresRepository) Update(entity blog_models.Blog) error {
	return b.base.Update(entity)
}

func (b *BlogPostgresRepository) Delete(id string) error {
	return b.base.Delete(id)
}

func (repo *BlogPostgresRepository) ListPublished() ([]blog_models.Blog, error) {
	list := []blog_models.Blog{}
	allPosts, _ := repo.List()
	for _, post := range allPosts {
		if post.IsPublished() {
			list = append(list, post)
		}
	}
	return list, nil
}

func CreatePostgresBlogRepo() *BlogPostgresRepository {

	repo := &BlogPostgresRepository{base: *db.CreatePostgresRepo[blog_models.Blog]("blogposts")}
	return repo
}
