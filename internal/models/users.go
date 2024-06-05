package models

import (
	"database/sql"
	"errors"
	"fmt"
	queryHelpers "vulh/soundcommunity/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int            `json:"id"`
	Username       string         `json:"username"`
	HashedPassword string         `json:"-"`
	Email          string         `json:"email"`
	CreatedAt      string         `json:"-"`
	UpdatedAt      sql.NullString `json:"-"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(createUserInput *CreateUserInput) error {
	stmt := `INSERT INTO users(username, hashed_password, email) VALUES($1, $2, $3)`
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(createUserInput.Password), 10)
	_, err := m.DB.Exec(stmt, createUserInput.Username, hashedPassword, createUserInput.Email)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetUserById(id int) (*User, error) {
	user := &User{}
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	row := queryBuilder.
		Select("id", "username", "email").
		FromTable("users").
		WhereColumn("id").
		Equal(fmt.Sprintf("%v", id)).
		GetOne()
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	return user, nil
}

func (m *UserModel) Login(loginInput *LoginInput) (*User, error) {
	user := &User{}
	queryBuilder := &queryHelpers.QueryBuilder{DB: m.DB}
	row := queryBuilder.
		Select("id", "username", "hashed_password", "email").
		FromTable("users").
		WhereColumn("username").
		Equal(loginInput.Username).
		GetOne()
	err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Email)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginInput.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}
