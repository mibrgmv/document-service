package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/mibrgmv/document-service/internal/repository"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	sql := `
	insert into users (id, login, password, created)
	values ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(ctx, sql, user.ID, user.Login, user.Password, user.Created)
	return err
}

func (r *userRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	sql := `
	select id, login, password, created
	from users 
	where login = $1
	`

	row := r.pool.QueryRow(ctx, sql, login)

	var user domain.User
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Created)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UserExists(ctx context.Context, login string) (bool, error) {
	sql := `
	select exists(select 1 from users where login = $1)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, sql, login).Scan(&exists)
	return exists, err
}
