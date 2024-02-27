package db_test

import (
	"ioprodz/common/db"
	"ioprodz/common/policies"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEntity struct {
	Id   string
	Data string
}

func (b TestEntity) GetId() string {
	return b.Id
}

func TestBaseRepository(t *testing.T) {

	db.NewMongoConnection()

	repos := []policies.BaseRepository[TestEntity]{db.CreateMongoRepo[TestEntity]("test_collection"), db.CreateMemoryRepo[TestEntity]()}

	for _, repo := range repos {
		repo.Create(TestEntity{Id: "test-id1", Data: "data1"})
		repo.Create(TestEntity{Id: "test-id2", Data: "data2"})

		t.Run("return inserted documents list", func(t *testing.T) {

			list, _ := repo.List()
			assert.Equal(t, "data1", list[0].Data)
			assert.Equal(t, "data2", list[1].Data)
		})

		t.Run("update inserted document by id", func(t *testing.T) {

			repo.Update(TestEntity{Id: "test-id1", Data: "data0"})

			blog, _ := repo.Get("test-id1")
			assert.Equal(t, "data0", blog.Data)
		})
	}

}
