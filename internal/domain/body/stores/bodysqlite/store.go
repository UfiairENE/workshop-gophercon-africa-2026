package bodysqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"uuid"
	"workshop/internal/domain/body"
	"workshop/internal/infra/db/sqlc"
)

// Store persists bodies using sqlc-generated SQLite queries.
type Store struct {
	conn    *sql.DB
	queries *db.Queries
}

var _ body.Store = (*Store)(nil)

// New constructs a body SQLite store over conn.
func New(conn *sql.DB) *Store {
	return &Store{conn: conn, queries: db.New(conn)}
}

// Find returns the body for id, or body.ErrNotFound.
func (s Store) Find(ctx context.Context, id uuid.UUID) (body.Body, error) {
	row, err := s.queries.BodyFind(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return body.Body{}, body.ErrNotFound
	} else if err != nil {
		return body.Body{}, err
	}
	return transform(row), nil
}

// FindByName returns the body called name, ignoring case, or body.ErrNotFound.
func (s Store) FindByName(ctx context.Context, name string) (body.Body, error) {
	row, err := s.queries.BodyFindByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return body.Body{}, body.ErrNotFound
	} else if err != nil {
		return body.Body{}, err
	}
	return transform(row), nil
}

// List returns the bodies matching filter.
func (s Store) List(ctx context.Context, filter body.ListFilter) ([]body.Body, error) {
	kinds, err := json.Marshal(filter.Kinds)
	if err != nil {
		return nil, err
	}

	// SQLite reads a negative limit as no limit.
	limit := int64(filter.Limit)
	if limit <= 0 {
		limit = -1
	}

	rows, err := s.queries.BodyList(ctx, db.BodyListParams{
		AnyKind: len(filter.Kinds) == 0,
		Kinds:   string(kinds),
		System:  filter.System,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}

	result := make([]body.Body, len(rows))
	for i, row := range rows {
		result[i] = transform(row)
	}
	return result, nil
}

// CreateMany saves bodies in a single transaction, skipping any whose
// key is already stored. One transaction rather than one per row is the
// difference between a moment and a minute for the full catalogue.
func (s Store) CreateMany(ctx context.Context, in []body.CreateParams) (int, error) {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	q := s.queries.WithTx(tx)
	added := 0
	for _, p := range in {
		n, err := q.BodyCreate(ctx, toCreateParams(uuid.New(), p))
		if err != nil {
			return 0, fmt.Errorf("store: creating %q: %w", p.Key, err)
		}
		added += int(n)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return added, nil
}
