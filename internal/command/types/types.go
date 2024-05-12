package types

import (
	"context"
	"database/sql"
	"os/exec"
	"time"
)

type Command struct {
	ID        int
	Raw       string
	Name      string
	ErrorMsg  sql.NullString
	Logs      sql.NullString
	Status    sql.NullString
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type CommandExec struct {
	Command
	Exec *exec.Cmd
}

type CommandService interface {
	Create(ctx context.Context, name, raw string) (Command, error)
	List(ctx context.Context) ([]Command, error)
	Get(ctx context.Context, id int) (Command, error)
	Stop(ctx context.Context, id int) error
}

type CommandRepo interface {
	List(ctx context.Context) (cmds []Command, err error)
	Get(ctx context.Context, id int) (Command, error)
	Create(ctx context.Context, name, raw string) (Command, error)
	Remove(ctx context.Context, id int) error
	SetSuccess(id int) error
	SetError(id int, error_msg string) error
	StdoutWriter(id int, bb []byte) error
}

type CreateCommandRequest struct {
	Name string `json:"name" validate:"required"`
	Raw  string `json:"raw" validate:"required"`
}

type CreateCommandResponse struct {
	ID int `json:"id"`
}

type GetCommandResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Raw       string    `json:"raw"`
	ErrorMsg  string    `json:"errorMsg"`
	Logs      string    `json:"logs"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
