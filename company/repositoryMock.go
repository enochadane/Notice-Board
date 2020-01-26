package company

import (
	"github.com/amthesonofGod/Notice-Board/entity"
)

var b = [2]entity.Company{
	entity.Company{Name: "moti", Email: "king@zmountain", Password: "12"},
	entity.Company{Name: "king", Email: "king@zvalley", Password: "123"},
}
var a []entity.Company = b[0:len(b)]

//CompanyRepositoryMock ...
type CompanyRepositoryMock interface {
	CompaniesMock() ([]entity.Company, []error)
	CompanyMock(id uint) (*entity.Company, []error)
	UpdateCompanyMock(company *entity.Company) (*entity.Company, []error)
	DeleteCompanyMock(id uint) (*entity.Company, []error)
	StoreCompanyMock(company *entity.Company) (*entity.Company, []error)
	// StoreSession(session *entity.CompanySession) (*entity.CompanySession, []error)
	// Session(uuid string) (*entity.CompanySession, []error)
	// DeleteSession(uuid string) (*entity.CompanySession, []error)
}

//CompaniesMock ...
func CompaniesMock() ([]entity.Company, error) {
	return a, nil
}

// CompanyMock ...
func CompanyMock(id uint) (*entity.Company, error) {

	for _, s := range a {
		if s.ID == id {
			return &s, nil
		}
	}

	return nil, nil
}

// UpdateCompanyMock ...
func UpdateCompanyMock(user *entity.Company) (*entity.Company, error) {
	a[1] = *user

	return user, nil

}

// DeleteCompanyMock ...
func DeleteCompanyMock(id uint) (*entity.Company, error) {
	for i, q := range a {
		if id == q.ID && i == 0 {

			c := a[i+1 : len(a)-1]
			for j := i; j < len(c); j++ {
				a[j] = c[j-i]
			}
			return &a[i], nil
		}
	}
	return nil, nil

}

//StoreCompanyMock ...
func StoreCompanyMock(user *entity.Company) (*entity.Company, error) {
	a[len(a)] = *user

	return user, nil
}
