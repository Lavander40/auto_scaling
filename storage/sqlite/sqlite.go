package sqlite

import (
	"auto_scaling/lib/e"
	"auto_scaling/storage"
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db  *sql.DB
	ctx context.Context
}

func New(ctx context.Context, path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.WrapErr("can't connect to db", err)
	}

	if err = db.Ping(); err != nil {
		return nil, e.WrapErr("can't connect to db", err)
	}

	return &Storage{db: db, ctx: ctx}, nil
}

func (s Storage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS calls (type INT, amount INT, user_name TEXT, created_at TIMESTAMP)`
	_, err := s.db.ExecContext(s.ctx, q)
	if err != nil {
		return e.WrapErr("can't create table", err)
	}

	return nil
}

func (s Storage) Save(call *storage.Call) error {
	q := `INSERT INTO calls (type, amount, user_name, created_at) VALUES (?, ?, ?, ?)`

	_, err := s.db.ExecContext(s.ctx, q, call.Type, call.Amount, call.UserName, call.CreatedAt)
	if err != nil {
		return e.WrapErr("can't save", err)
	}

	return nil
}

func (s Storage) PickLastCalls(userName string) ([]*storage.Call, error) {
	q := `SELECT type, amount, user_name, created_at FROM calls WHERE user_name = ? ORDER BY created_at DESC LIMIT 10`

	rows, err := s.db.QueryContext(s.ctx, q, userName)
	if err != nil {
		return nil, e.WrapErr("can't save", err)
	}
	defer rows.Close()

	var calls []*storage.Call
	for rows.Next() {
		var call storage.Call
		if err := rows.Scan(&call.Type, &call.Amount, &call.UserName, &call.CreatedAt); err != nil {
            return nil, e.WrapErr("can't scan call", err)
        }
        calls = append(calls, &call)
	}
	if err := rows.Err(); err != nil {
        return nil, e.WrapErr("error iterating over rows", err)
    }
	if len(calls) == 0 {
		return nil, storage.ErrEmpty
	}

	return calls, nil
}
