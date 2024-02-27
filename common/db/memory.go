package db

import (
	"ioprodz/common/policies"
)

type BaseMemoryRepository[T policies.Entity] struct {
	list []T
}

func (repo *BaseMemoryRepository[T]) List() ([]T, error) {
	return repo.list, nil
}

func (repo *BaseMemoryRepository[T]) Create(entity T) error {
	repo.list = append(repo.list, entity)
	return nil
}

func (repo *BaseMemoryRepository[T]) Get(id string) (T, error) {
	var result T
	for _, obj := range repo.list {
		if obj.GetId() == id {
			result = obj
		}
	}
	return result, &policies.StorageError{Message: "Element not found by id: " + id}
}

func (repo *BaseMemoryRepository[T]) Update(entity T) error {
	for index, existing := range repo.list {
		if existing.GetId() == entity.GetId() {
			repo.list[index] = entity
			return nil
		}
	}
	return &policies.StorageError{Message: "Element not found by id: " + entity.GetId()}
}

func (repo *BaseMemoryRepository[T]) Delete(id string) error {

	return &policies.StorageError{Message: "Element not found by id: "}
}

func CreateMemoryRepo[T policies.Entity]() *BaseMemoryRepository[T] {
	repo := &BaseMemoryRepository[T]{list: []T{}}
	return repo
}