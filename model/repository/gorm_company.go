package repository

import (
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/model"
	"github.com/jinzhu/gorm"
)

// CompanyGormRepo implements the model.CompanyRepository interface
type CompanyGormRepo struct {
	conn *gorm.DB
}

// NewCompanyGormRepo will create a new object of CompanyGormRepo
func NewCompanyGormRepo(db *gorm.DB) model.CompanyRepository {
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

// StoreSession stores a given session in the database
func (cRepo *CompanyGormRepo) StoreSession(session *entity.Session) (*entity.Session, []error) {
	s := session
	errs := cRepo.conn.Create(s).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return s, errs
}

// DeleteSession deletes a given session from the database
func (cRepo *CompanyGormRepo) DeleteSession(uuid string) (*entity.Session, []error) {
	s, errs := cRepo.Session(uuid)
	if len(errs) > 0 {
		return nil, errs
	}

	errs = cRepo.conn.Delete(s, s.UUID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return s, errs
}

// Session retrieve a session from the database by its id
func (cRepo *CompanyGormRepo) Session(uuid string) (*entity.Session, []error) {
	s := entity.Session{}
	errs := cRepo.conn.Where("UUID = ?", uuid).First(&s).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &s, errs
}
