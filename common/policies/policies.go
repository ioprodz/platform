package policies

type Entity interface {
	GetId() string
}

type BaseRepository[T Entity] interface {
	List() ([]T, error)
	Get(id string) (T, error)
	Delete(id string) error
	Create(entity T) error
	Update(entity T) error
}

type StorageError struct {
	Message string
}

func (e *StorageError) Error() string {
	return e.Message
}
