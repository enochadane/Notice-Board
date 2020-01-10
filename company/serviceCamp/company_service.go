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

// StoreSession persists new session information
func (cs *CompanyService) StoreSession(session *entity.CompanySession) (*entity.CompanySession, []error) {

	s, errs := cs.companyRepo.StoreSession(session)

	if len(errs) > 0 {
		return nil, errs
	}

	return s, nil
}

// DeleteSession delete a session by its id
func (cs *CompanyService) DeleteSession(uuid string) (*entity.CompanySession, []error) {

	s, errs := cs.companyRepo.DeleteSession(uuid)

	if len(errs) > 0 {
		return nil, errs
	}
	return s, nil
}

// Session returns a session object with a given id
func (cs *CompanyService) Session(uuid string) (*entity.CompanySession, []error) {

	s, errs := cs.companyRepo.Session(uuid)

	if len(errs) > 0 {
		return nil, errs
	}

	return s, nil
}
