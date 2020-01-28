package repository

import (
	"errors"

	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/jinzhu/gorm"
)

// CompanyGormRepoMock implements the model.CompanyRepository interface
type CompanyGormRepoMock struct {
	conn *gorm.DB
}

// NewCompanyGormRepoMock will create a new object of CompanyGormRepo
func NewCompanyGormRepoMock(db *gorm.DB) company.CompanyRepository {
	return &CompanyGormRepoMock{conn: db}
}

// Companies returns all companies stored in the database
func (cRepo *CompanyGormRepoMock) Companies() ([]entity.Company, []error) {
	cmps := []entity.Company{entity.CompanyMock}

	return cmps, nil
}

// Company retrieve a company from the database by its id
func (cRepo *CompanyGormRepoMock) Company(id uint) (*entity.Company, []error) {

	cmp := entity.CompanyMock
	if id == 1 {
		return &cmp, nil

	}
	return nil, []error{errors.New("post not found")}
}

// UpdateCompany updates a given company in the database
func (cRepo *CompanyGormRepoMock) UpdateCompany(company *entity.Company) (*entity.Company, []error) {
	cmp := entity.CompanyMock

	return &cmp, nil
}

// DeleteCompany deletes a given company from the database
func (cRepo *CompanyGormRepoMock) DeleteCompany(id uint) (*entity.Company, []error) {
	cmp := entity.CompanyMock
	if id != 1 {
		return nil, []error{errors.New("post not found")}
	}
	return &cmp, nil

}

// StoreCompany stores a given company in the database
func (cRepo *CompanyGormRepoMock) StoreCompany(company *entity.Company) (*entity.Company, []error) {
	cmp := company

	return cmp, nil
}

func (cRepo *CompanyGormRepoMock) CompanyByEmail(email string) (*entity.Company, []error) {
	company := entity.Company{}
	errs := cRepo.conn.Find(&company, "email=?", email).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &company, errs
}

// PhoneExists check if a given phone number is found
func (cRepo *CompanyGormRepoMock) PhoneExists(phone string) bool {

	return true
}

// EmailExists check if a given email is found
func (cRepo *CompanyGormRepoMock) EmailExists(email string) bool {
	comp := entity.CompanyMock

	if comp.Email != email {
		return false
	}
	return true
}

//CompanyRoles returns list of application roles that a given user has
func (cRepo *CompanyGormRepoMock) CompanyRoles(company *entity.Company) ([]entity.Role, []error) {
	companyRoles := []entity.Role{entity.Rol}
	return companyRoles, nil

}

// // StoreSession stores a given session in the database
// func (cRepo *CompanyGormRepoMock) StoreSession(session *entity.CompanySession) (*entity.CompanySession, []error) {
// 	s := session

// 	return s, nil
// }

// // DeleteSession deletes a given session from the database
// func (cRepo *CompanyGormRepoMock) DeleteSession(uuid string) (*entity.CompanySession, []error) {
// 	s := entity.CampSessionMock

// 	return &s, nil
// }

// // Session retrieve a session from the database by its id
// func (cRepo *CompanyGormRepoMock) Session(uuid string) (*entity.CompanySession, []error) {
// 	s := entity.CampSessionMock

// 	return &s, []error{errors.New("post not found")}
// }
