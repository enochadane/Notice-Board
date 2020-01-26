package repository

import (
	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/jinzhu/gorm"
)

// CompanyGormRepo implements the model.CompanyRepository interface
type CompanyGormRepo struct {
	conn *gorm.DB
}

// NewCompanyGormRepo will create a new object of CompanyGormRepo
func NewCompanyGormRepo(db *gorm.DB) company.CompanyRepository {
	return &CompanyGormRepo{conn: db}
}

// Companies returns all companies stored in the database
func (cRepo *CompanyGormRepo) Companies() ([]entity.Company, []error) {
	cmps := []entity.Company{}
	errs := cRepo.conn.Find(&cmps).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cmps, errs
}

// Company retrieve a company from the database by its id
func (cRepo *CompanyGormRepo) Company(id uint) (*entity.Company, []error) {
	cmp := entity.Company{}
	errs := cRepo.conn.First(&cmp, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &cmp, errs
}

// UpdateCompany updates a given company in the database
func (cRepo *CompanyGormRepo) UpdateCompany(company *entity.Company) (*entity.Company, []error) {
	cmp := company
	errs := cRepo.conn.Save(cmp).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cmp, errs
}

// DeleteCompany deletes a given company from the database
func (cRepo *CompanyGormRepo) DeleteCompany(id uint) (*entity.Company, []error) {
	cmp, errs := cRepo.Company(id)
	if len(errs) > 0 {
		return nil, errs
	}

	errs = cRepo.conn.Delete(cmp, cmp.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cmp, errs
}

// StoreCompany stores a given company in the database
func (cRepo *CompanyGormRepo) StoreCompany(company *entity.Company) (*entity.Company, []error) {
	cmp := company
	errs := cRepo.conn.Create(cmp).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cmp, errs
}

// CompanyByEmail retrieves a company user by its email address from the database
func (cRepo *CompanyGormRepo) CompanyByEmail(email string) (*entity.Company, []error) {
	company := entity.Company{}
	errs := cRepo.conn.Find(&company, "email=?", email).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &company, errs
}

// PhoneExists check if a given phone number is found
func (cRepo *CompanyGormRepo) PhoneExists(phone string) bool {
	comp := entity.Company{}
	errs := cRepo.conn.Find(&comp, "phone=?", phone).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}

// EmailExists check if a given email is found
func (cRepo *CompanyGormRepo) EmailExists(email string) bool {
	comp := entity.Company{}
	errs := cRepo.conn.Find(&comp, "email=?", email).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}

// UserRoles returns list of application roles that a given user has
// func (userRepo *UserGormRepo) UserRoles(user *entity.User) ([]entity.Role, []error) {
// 	userRoles := []entity.Role{}
// 	errs := userRepo.conn.Model(user).Related(&userRoles).GetErrors()
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return userRoles, errs
// }
