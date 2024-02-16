package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nalivayko13/codingTask/store/model"
)

func NewPostgresRepo(db *sql.DB) *RepoPostgres {
	return &RepoPostgres{db: db}
}

func (r *RepoPostgres) CreateStore(ctx context.Context, store *model.Store) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		tx.Rollback()
		return err
	}

	//create new store
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO "stores" (name, address, created_at)
 			VALUES ($1, $2, $3) RETURNING id`,
		store.Name, store.Address, store.CreatedAt)
	var storeID int
	if err := row.Scan(&storeID); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *RepoPostgres) GetStore(ctx context.Context, id int) (*model.Store, error) {
	var store model.Store

	result := r.db.QueryRowContext(ctx,
		`SELECT id, name, address, created_at FROM stores WHERE id = $1 and is_deleted=false`, id)
	if err := result.Scan(&store.ID, &store.Name, &store.Address, &store.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("store not found: %w", err)
		}
		return nil, err
	}
	return &store, nil
}

func (r *RepoPostgres) GetVersions(ctx context.Context, storeID int) ([]model.Version, error) {
	var versions []model.Version

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, store_id, version_number, creator, owner, open_at, close_at, created_at
				FROM versions WHERE store_id = $1 ORDER BY created_at DESC`, storeID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var version model.Version
		if err := rows.Scan(&version.ID, &version.StoreID, &version.VersionNumber, &version.Creator,
			&version.Owner, &version.OpenAt, &version.CloseAt, &version.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (r *RepoPostgres) GetVersion(ctx context.Context, storeID, versionID int) (*model.Version, error) {
	var version model.Version

	result := r.db.QueryRowContext(ctx,
		`SELECT id, store_id, version_number, creator, owner, created_at, open_at, close_at
					FROM versions WHERE id = $1 and store_id=$2 ORDER BY created_at DESC`, versionID, storeID)
	if err := result.Scan(&version.ID, &version.StoreID, &version.VersionNumber, &version.Creator,
		&version.Owner, &version.CreatedAt, &version.OpenAt, &version.CloseAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("version not found: %w", err)
		}
		return nil, err
	}
	return &version, nil
}

func (r *RepoPostgres) CreateVersion(ctx context.Context, version *model.Version) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		tx.Rollback()
		return err
	}

	//create new store
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO "versions" (store_id, version_number, creator, owner, open_at, close_at, created_at)
 			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		version.StoreID, version.VersionNumber, version.Creator, version.Owner,
		version.OpenAt, version.CloseAt, version.CreatedAt)
	var storeID int
	if err := row.Scan(&storeID); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *RepoPostgres) DeleteStore(ctx context.Context, storeID int) error {
	var ID int
	row := r.db.QueryRowContext(ctx,
		`UPDATE "stores" SET is_deleted = 'true' WHERE id = $1 RETURNING id`, storeID)
	if err := row.Scan(&ID); err != nil {
		return err
	}
	return nil
}

func (r *RepoPostgres) DeleteVersion(ctx context.Context, versionID, storeID int) error {
	var ID int
	row := r.db.QueryRowContext(ctx,
		`UPDATE "versions" SET is_deleted = 'true' WHERE id = $1 and store_id = $2 RETURNING id`,
		versionID, storeID)
	if err := row.Scan(&ID); err != nil {
		return err
	}
	return nil
}
