package repository

import (
	"database/sql"
	"errors"

	"github.com/amthesonofGod/Notice-Board/entity"
)

type CompanyRepositoryImpl struct {
	conn *sql.DB
}

func NewCompanyRepositoryImpl(Conn *sql.DB) *CompanyRepositoryImpl {
	return &CompanyRepositoryImpl{conn: Conn}
}

func (ci *CompanyRepositoryImpl) Companies() ([]entity.Company, error) {

	rows, err := ci.conn.Query("SELECT * FROM companies;")
	if err != nil {
		return nil, errors.New("Could not query the database")
	}
	defer rows.Close()

	cmpns := []entity.Company{}

	for rows.Next() {
		cmp := entity.Company{}
		err = rows.Scan(&cmp.ID, &cmp.Name, &cmp.Email, &cmp.Password)
		if err != nil {
			return nil, err
		}
		cmpns = append(cmpns, cmp)
	}

	return cmpns, nil

}

func (ci *CompanyRepositoryImpl) Company(id int) (entity.Company, error) {

	row := ci.conn.QueryRow("SELECT * FROM companies WHERE id = $1", id)

	c := entity.Company{}

	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Password)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (ci *CompanyRepositoryImpl) UpdateCompany(c entity.Company) error {

	_, err := ci.conn.Exec("UPDATE companies SET name=$1,email=$2, password=$3 WHERE id=$4", c.Name, c.Email, c.Password, c.ID)
	if err != nil {
		return errors.New("Update has failed")
	}

	return nil
}

func (ci *CompanyRepositoryImpl) DeleteCompany(id int) error {

	_, err := ci.conn.Exec("DELETE FROM companies WHERE id=$1", id)
	if err != nil {
		return errors.New("Delete has failed")
	}

	return nil
}

func (ci *CompanyRepositoryImpl) StoreCompany(c entity.Company) error {

	_, err := ci.conn.Exec("INSERT INTO companies (name,email,password) values($1, $2, $3)", c.Name, c.Email, c.Password)
	if err != nil {
		return errors.New("Insertion has failed")
	}

	return nil
}
