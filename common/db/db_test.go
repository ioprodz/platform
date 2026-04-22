package db_test

import (
	"context"
	"ioprodz/common/db"
	"ioprodz/common/policies"
	"os"
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

	repos := []policies.BaseRepository[TestEntity]{db.CreateMemoryRepo[TestEntity]()}

	if os.Getenv("DATABASE_URL") != "" {
		pool := db.GetPool()
		_, err := pool.Exec(context.Background(),
			"CREATE TABLE IF NOT EXISTS test_entities (id TEXT PRIMARY KEY, data JSONB NOT NULL)")
		if err != nil {
			t.Fatalf("failed to create test_entities table: %v", err)
		}
		_, _ = pool.Exec(context.Background(), "TRUNCATE test_entities")
		t.Cleanup(func() {
			_, _ = pool.Exec(context.Background(), "DROP TABLE test_entities")
		})
		repos = append(repos, db.CreatePostgresRepo[TestEntity]("test_entities"))
	}

	for _, repo := range repos {
		repo.Create(TestEntity{Id: "test-id1", Data: "data1"})
		repo.Create(TestEntity{Id: "test-id2", Data: "data2"})

		t.Run("return inserted documents list", func(t *testing.T) {

			list, _ := repo.List()
			assert.Equal(t, 2, len(list))
			ids := map[string]string{}
			for _, e := range list {
				ids[e.Id] = e.Data
			}
			assert.Equal(t, "data1", ids["test-id1"])
			assert.Equal(t, "data2", ids["test-id2"])
		})

		t.Run("update inserted document by id", func(t *testing.T) {

			repo.Update(TestEntity{Id: "test-id1", Data: "data0"})

			entity, _ := repo.Get("test-id1")
			assert.Equal(t, "data0", entity.Data)
		})
	}

}
