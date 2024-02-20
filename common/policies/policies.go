package policies

type Entity interface {
}

type Repository[T Entity] interface {
	List() []T
	Get(id string) (T, error)
	Create(entity T)
	// Update(entity T)
	// Delete(entity T)
}
