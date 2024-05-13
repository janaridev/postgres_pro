package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/janaridev/postgres_pro/internal/command/repo"
	"github.com/janaridev/postgres_pro/internal/command/types"
	"github.com/janaridev/postgres_pro/pkg/logger/sl"
	"github.com/janaridev/postgres_pro/pkg/pool"
	"github.com/janaridev/postgres_pro/pkg/xmap"
)

type CommandService struct {
	log     *slog.Logger
	repo    types.CommandRepo
	pool    *pool.Pool
	workers xmap.XMap[int, *types.CommandExec]
}

func New(log *slog.Logger, repo types.CommandRepo, pool *pool.Pool) *CommandService {
	return &CommandService{
		log:  log,
		repo: repo,
		pool: pool,
	}
}

var (
	ErrCommandNotFound      = errors.New("command not found")
	ErrCommandAlreadyExists = errors.New("command already exists")
)

func (c *CommandService) List(ctx context.Context) ([]types.Command, error) {
	const op = "commandService.List"

	cmds, err := c.repo.List(ctx)
	if err != nil {
		c.log.Warn(err.Error())

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cmds, nil
}

func (c *CommandService) Get(ctx context.Context, id int) (types.Command, error) {
	const op = "commandService.Get"

	cmd, err := c.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrCommandNotFound) {
			c.log.Warn(err.Error())

			return types.Command{}, fmt.Errorf("%s: %w", op, ErrCommandNotFound)
		}

		c.log.Warn(err.Error())

		return types.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	return cmd, err
}

func (c *CommandService) Create(ctx context.Context, name, raw string) (types.Command, error) {
	const op = "commandService.Create"

	logger := c.log.With(slog.String("op", op))

	cmd, err := c.repo.Create(ctx, name, raw)
	if err != nil {
		if errors.Is(err, repo.ErrCommandExists) {
			return types.Command{}, fmt.Errorf("%s: %w", op, ErrCommandAlreadyExists)
		}

		return types.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	c.pool.Go(func() {
		if err := c.exec(cmd); err != nil {
			logger.ErrorContext(ctx, "command exec", sl.Err(err))
		}
	})

	return cmd, nil
}

func (c *CommandService) Stop(ctx context.Context, id int) error {
	defer func() {
		c.workers.Delete(id)
	}()

	load, ok := c.workers.Load(id)
	if !ok {
		return ErrCommandNotFound
	}

	if err := c.repo.Remove(ctx, id); err != nil {
		c.log.Warn(err.Error())

		return err
	}

	return load.Exec.Process.Kill()
}

func (c *CommandService) Launch(ctx context.Context, id int) error {
	const op = "commandService.Launch"

	var cmd types.Command

	load, ok := c.workers.Load(id)
	if ok {
		cmd = load.Command
	} else {
		res, err := c.repo.Get(ctx, id)
		if err != nil {
			if errors.Is(err, repo.ErrCommandNotFound) {
				c.log.Warn(err.Error())

				return ErrCommandNotFound
			}

			c.log.Warn(err.Error())
			return err
		}

		cmd = res
	}

	c.pool.Go(func() {
		if err := c.exec(cmd); err != nil {
			c.log.ErrorContext(ctx, fmt.Sprintf("%s command exec", op), sl.Err(err))
		}
	})

	return nil
}

func (c *CommandService) exec(cmd types.Command) (err error) {
	const op = "commandService.exec"

	logger := c.log.With(slog.String("op", op))

	defer func() {
		if err == nil {
			err = c.repo.SetSuccess(cmd.ID)
			return
		}
		logger.Warn("something went wrong while executing command", sl.Err(err))

		err = errors.Join(err, c.repo.SetError(cmd.ID, err.Error()))
	}()

	execCmd := exec.Command("/bin/sh", "-c", cmd.Raw)
	stdoutPipe, err := execCmd.StdoutPipe()
	if err != nil {
		logger.Warn("failed to create stdout pipe", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	c.workers.Store(cmd.ID, &types.CommandExec{
		Command: cmd,
		Exec:    execCmd,
	})

	if err := execCmd.Start(); err != nil {
		logger.Warn("failed to start command", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdoutPipe.Read(buf)
			if err != nil {
				break
			}

			if err := c.repo.StdoutWriter(cmd.ID, buf[:n]); err != nil {
				break
			}

			logger.Info("created log")
		}
	}()

	if err = execCmd.Wait(); err != nil {
		logger.Warn("failed to wait command", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
