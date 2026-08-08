package postgres

import (
	"context"
	"database/sql"
	"errors"

	"example.com/htmahs/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

// ensure this implements all the methods in user.Repository
var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	const query = `
		INSERT INTO users
			(name, email)
		VALUES
			($1, $2)
		RETURNING 
			id, created_at, updated_at;
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	const query = `
		SELECT 
			id, name, email, created_at, updated_at
		FROM 
			users 
		WHERE 
			id = $1;
	`

	u := &user.User{}

	err := r.db.QueryRowContext(
		ctx, query, id,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, id string) (*user.User, error) {
	const query = `
		SELECT 
			id, name, email, created_at, updated_at
		FROM 
			users 
		WHERE 
			id = $1;
	`

	u := &user.User{}

	err := r.db.QueryRowContext(
		ctx, query, id,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]user.User, error) {
	query := `
		SELECT 
			id, name, email, created_at, updated_at
		FROM
			users
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []user.User

	for rows.Next() {
		var u user.User

		rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt,
		)

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
