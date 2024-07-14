package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/phonghaido/artifactory-upload-service/helpers"
	"github.com/phonghaido/artifactory-upload-service/types"
)

type PostgreSQL struct {
	Conn string
}

func NewPostgreSQL(conn string) *PostgreSQL {
	return &PostgreSQL{
		Conn: conn,
	}
}

func (p *PostgreSQL) SearchUser(username, password string) (types.User, error) {
	db, err := sql.Open("postgres", p.Conn)
	if err != nil {
		return types.User{}, err
	}
	defer db.Close()

	query := "SELECT * FROM \"user\" WHERE username = $1 AND password = $2"
	row := db.QueryRow(query, username, password)

	var user types.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Scope); err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, helpers.Unauthorized()
		}
		return types.User{}, err
	}

	return user, nil
}
