package db

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"ioprodz/common/config"
	"ioprodz/common/policies"
	"log"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var pool *pgxpool.Pool

func GetPool() *pgxpool.Pool {
	if pool != nil {
		return pool
	}
	uri := config.Load().DATABASE_URL
	if uri == "" {
		panic("DATABASE_URL missing")
	}
	p, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		panic(err)
	}
	if err := p.Ping(context.Background()); err != nil {
		panic(err)
	}
	pool = p
	return pool
}

func RunMigrations() {
	uri := config.Load().DATABASE_URL
	if uri == "" {
		panic("DATABASE_URL missing")
	}

	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(err)
	}
	src, err := iofs.New(sub, ".")
	if err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(GetPool())
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "pgx5", driver)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	log.Println("db: migrations applied")
}

type BasePostgresRepository[T policies.Entity] struct {
	tableName string
}

func (repo *BasePostgresRepository[T]) List() ([]T, error) {
	list := []T{}
	rows, err := GetPool().Query(context.Background(), "SELECT data FROM "+repo.tableName)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return list, err
		}
		var t T
		if err := json.Unmarshal(raw, &t); err != nil {
			return list, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (repo *BasePostgresRepository[T]) Get(id string) (T, error) {
	var raw []byte
	err := GetPool().QueryRow(context.Background(),
		"SELECT data FROM "+repo.tableName+" WHERE id = $1", id).Scan(&raw)
	if err != nil {
		var empty T
		return empty, &policies.StorageError{Message: "Element not found by id: " + id + " (" + err.Error() + ")"}
	}
	var t T
	if err := json.Unmarshal(raw, &t); err != nil {
		var empty T
		return empty, err
	}
	return t, nil
}

func (repo *BasePostgresRepository[T]) Create(entity T) error {
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	_, err = GetPool().Exec(context.Background(),
		"INSERT INTO "+repo.tableName+" (id, data) VALUES ($1, $2)", entity.GetId(), data)
	return err
}

func (repo *BasePostgresRepository[T]) Update(entity T) error {
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	tag, err := GetPool().Exec(context.Background(),
		"UPDATE "+repo.tableName+" SET data = $1 WHERE id = $2", data, entity.GetId())
	if err != nil {
		return &policies.StorageError{Message: "Element not found by id: " + entity.GetId()}
	}
	if tag.RowsAffected() == 0 {
		return &policies.StorageError{Message: "Element not found by id: " + entity.GetId()}
	}
	return nil
}

func (repo *BasePostgresRepository[T]) Delete(id string) error {
	_, err := GetPool().Exec(context.Background(),
		"DELETE FROM "+repo.tableName+" WHERE id = $1", id)
	if err != nil {
		return &policies.StorageError{Message: "Element could not be deleted by id: " + id}
	}
	return nil
}

func CreatePostgresRepo[T policies.Entity](tableName string) *BasePostgresRepository[T] {
	return &BasePostgresRepository[T]{tableName: tableName}
}
