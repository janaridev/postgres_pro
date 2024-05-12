package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/janaridev/postgres_pro/internal/command/types"
)

type CommandRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *CommandRepo {
	return &CommandRepo{
		db: db,
	}
}

var (
	ErrCommandExists   = errors.New("command already exists")
	ErrCommandNotFound = errors.New("command not found")
)

func (c *CommandRepo) List(ctx context.Context) (cmds []types.Command, err error) {
	const op = "commandRepo.List"

	const query = `
		SELECT c.id, c.name, c.raw, c.created_at, c.updated_at, c.error_msg, c.status, cl.logs
		FROM commands as c
		LEFT JOIN command_logs as cl ON c.id = cl.command_id
	`

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	for rows.Next() {
		var cmd types.Command
		if err := rows.Scan(&cmd.ID, &cmd.Name, &cmd.Raw,
			&cmd.CreatedAt, &cmd.UpdatedAt,
			&cmd.ErrorMsg, &cmd.Status, &cmd.Logs); err != nil {
			if err == sql.ErrNoRows {
				return cmds, nil
			}

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		cmds = append(cmds, cmd)
	}

	return cmds, err
}

func (c *CommandRepo) Get(ctx context.Context, id int) (cmd types.Command, err error) {
	const op = "commandRepo.Get"

	const query = `
		SELECT c.id, c.name, c.raw, c.created_at, c.updated_at, c.error_msg, c.status, cl.logs
		FROM commands as c
		LEFT JOIN command_logs as cl ON c.id = cl.command_id
		WHERE c.id = $1
	`

	row := c.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return types.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := row.Scan(&cmd.ID, &cmd.Name, &cmd.Raw,
		&cmd.CreatedAt, &cmd.UpdatedAt,
		&cmd.ErrorMsg, &cmd.Status, &cmd.Logs); err != nil {
		if err == sql.ErrNoRows {
			return types.Command{}, fmt.Errorf("%s: command not found: %w", op, ErrCommandNotFound)
		}

		return types.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	return cmd, nil
}

func (c *CommandRepo) Create(ctx context.Context, name, raw string) (cmd types.Command, err error) {
	const op = "commandRepo.Create"

	const query = `
		INSERT INTO commands(name, raw) 
		VALUES ($1, $2) 
		ON CONFLICT (name) DO UPDATE
		SET raw = $2
		RETURNING id, raw, created_at
	`

	row := c.db.QueryRowContext(ctx, query, name, raw)
	if err := row.Err(); err != nil {
		return types.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := row.Scan(&cmd.ID, &cmd.Raw, &cmd.CreatedAt); err != nil {
		return types.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	return cmd, err
}

func (c *CommandRepo) Remove(ctx context.Context, id int) error {
	const op = "commandRepo.Remove"

	const query = `UPDATE commands SET is_deleted = true WHERE id = $1`
	
	if _, err := c.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *CommandRepo) SetSuccess(id int) error {
	const op = "commandRepo.SetSuccess"

	const query = `UPDATE commands SET status = 'success' WHERE id = $1`

	if _, err := c.db.Exec(query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *CommandRepo) SetError(id int, error_msg string) error {
	const op = "commandRepo.SetError"

	const query = `UPDATE commands SET status = 'error', error_msg = $1 WHERE id = $2`

	if _, err := c.db.Exec(query, error_msg, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *CommandRepo) StdoutWriter(id int, bb []byte) error {
	const op = "commandRepo.StdoutWriter"

	const query = `
		INSERT INTO command_logs(command_id, logs) 
		VALUES ($1, $2)
		ON CONFLICT (command_id) DO UPDATE
		SET logs = $2
	`

	if _, err := c.db.Exec(query, id, bb); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
