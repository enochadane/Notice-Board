package service

import (
	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
)

// CompanyService implements model.CompanyRepository interface
type CompanyService struct {
	companyRepo company.CompanyRepository
}

// NewCompanyService will create new CompanyService object
func NewCompanyService(CompanyRepos company.CompanyRepository) company.CompanyService {
	return &CompanyService{companyRepo: CompanyRepos}
}

// Companies returns list of companies
func (cs *CompanyService) Companies() ([]entity.Company, []error) {

	companies, errs := cs.companyRepo.Companies()

	if len(errs) > 0 {
		return nil, errs
	}

	return companies, nil
}

// StoreCompany persists new company information
func (cs *CompanyService) StoreCompany(company *entity.Company) (*entity.Company, []error) {

	cmp, errs := cs.companyRepo.StoreCompany(company)

	if len(errs) > 0 {
		return nil, errs
	}

	return cmp, nil
}

// Company returns a company object with a given id
func (cs *CompanyService) Company(id uint) (*entity.Company, []error) {

	cmp, errs := cs.companyRepo.Company(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return cmp, nil
}

// UpdateCompany updates a company with new data
func (cs *CompanyService) UpdateCompany(company *entity.Company) (*entity.Company, []error) {

	cmp, errs := cs.companyRepo.UpdateCompany(company)

	if len(errs) > 0 {
		return nil, errs
	}

	return cmp, nil
}

// DeleteCompany delete a company by its id
func (cs *CompanyService) DeleteCompany(id uint) (*entity.Company, []error) {

	cmp, errs := cs.companyRepo.DeleteCompany(id)

	if len(errs) > 0 {
		return nil, errs
	}
	return cmp, nil
}

// CompanyByEmail retrieves an application company user by its email address
func (cs *CompanyService) CompanyByEmail(email string) (*entity.Company, []error) {
	cmp, errs := cs.companyRepo.CompanyByEmail(email)
	if len(errs) > 0 {
		return nil, errs
	}
	return cmp, errs
}

// PhoneExists check if there is a user with a given phone number
func (cs *CompanyService) PhoneExists(phone string) bool {
	exists := cs.companyRepo.PhoneExists(phone)
	return exists
}

// EmailExists checks if there exist a user with a given email address
func (cs *CompanyService) EmailExists(email string) bool {
	exists := cs.companyRepo.EmailExists(email)
	return exists
}

// UserRoles returns list of roles a user has
func (cs *CompanyService) UserRoles(company *entity.Company) ([]entity.Role, []error) {
	companyRoles, errs := cs.companyRepo.UserRoles(company)
	if len(errs) > 0 {
		return nil, errs
	}
	return companyRoles, errs
}