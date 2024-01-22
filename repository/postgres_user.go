package repository

import (
	"context"

	"github.com/GermanPachec0/app-go/domain"
	"github.com/google/uuid"
)

type postgresUserRepository struct {
	conn Connection
}

func NewPostgresUser(conn Connection) domain.UserRepository {
	return &postgresUserRepository{conn}
}

func (p *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.User, error) {
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uu []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID,
			&u.Uuid,
			&u.Email,
			&u.Password,
			&u.FirstName,
			&u.LastName,
			&u.CreateAt,
			&u.UpdateAt,
		); err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}
	return uu, nil

}

func (p *postgresUserRepository) GetById(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	query := `SELECT id, uuid, email, password, first_name, last_name, created_at,update_at from
	users
	where uuid = $1`
	srs, err := p.fetch(ctx, query, uuid)
	if err != nil {
		return domain.User{}, err
	}
	if len(srs) == 0 {
		return domain.User{}, domain.ErrNotFound
	}
	return srs[0], nil
}
func (p *postgresUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	return nil, nil
}
func (p *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	query := `INSERT INTO users(email,password,first_name,last_name,create_at,update_at)
	VALUES($1,$2,NOW(),NOW())`

	return p.conn.QueryRow(ctx,
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName).Scan(&user.Uuid)
}

func (p *postgresUserRepository) Update(ctx context.Context) (*domain.User, error) {
	return nil, nil
}
