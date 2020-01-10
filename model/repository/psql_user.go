package repository

import (
	"database/sql"
	"errors"

	"github.com/amthesonofGod/Notice-Board/entity"
)

type UserRepositoryImpl struct {
	conn *sql.DB
}

func NewUserRepositoryImpl(Conn *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{conn: Conn}
}

func (usi *UserRepositoryImpl) Users() ([]entity.User, error) {

	rows, err := usi.conn.Query("SELECT * FROM users;")
	if err != nil {
		return nil, errors.New("Could not query the database")
	}
	defer rows.Close()

	usrs := []entity.User{}

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		usrs = append(usrs, user)
	}

	return usrs, nil

}

func (usi *UserRepositoryImpl) User(id int) (entity.User, error) {

	row := usi.conn.QueryRow("SELECT * FROM users WHERE id = $1", id)

	u := entity.User{}

	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (usi *UserRepositoryImpl) UpdateUser(u entity.User) error {

	_, err := usi.conn.Exec("UPDATE users SET name=$1,email=$2, password=$3 WHERE id=$4", u.Name, u.Email, u.Password, u.Id)
	if err != nil {
		return errors.New("Update has failed")
	}

	return nil
}

func (usi *UserRepositoryImpl) DeleteUser(id int) error {

	_, err := usi.conn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return errors.New("Delete has failed")
	}

	return nil
}

func (usi *UserRepositoryImpl) StoreUser(u entity.User) error {

	_, err := usi.conn.Exec("INSERT INTO users (name,email,password) values($1, $2, $3)", u.Name, u.Email, u.Password)
	if err != nil {
		return errors.New("Insertion has failed")
	}

	return nil
}