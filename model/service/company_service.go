package service

import (
	"github.com/motikingo/Notice-Board/entity"
	"github.com/motikingo/Notice-Board/model"
)

//CompanyServiceImpl ...
type CompanyServiceImpl struct {
	companyRepo model.CompanyRepository
}

// NewCompanyServiceImpl ...
func NewCompanyServiceImpl(CompanyRepos model.CompanyRepository) *CompanyServiceImpl {
	return &CompanyServiceImpl{companyRepo: CompanyRepos}
}

//Companies ...
func (cs *CompanyServiceImpl) Companies() ([]entity.Company, error) {

	companies, err := cs.companyRepo.Companies()

	if err != nil {
		return nil, err
	}

	return companies, nil
}

// StoreCompany ...
func (cs *CompanyServiceImpl) StoreCompany(company entity.Company) error {

	err := cs.companyRepo.StoreCompany(company)

	if err != nil {
		return err
	}

	return nil
}

// Company ...
func (cs *CompanyServiceImpl) Company(id int) (entity.Company, error) {

	c, err := cs.companyRepo.Company(id)

	if err != nil {
		return c, err
	}

	return c, nil
}

// UpdateCompany ...
func (cs *CompanyServiceImpl) UpdateCompany(company entity.Company) error {

	err := cs.companyRepo.UpdateCompany(company)

	if err != nil {
		return err
	}

	return nil
}

// DeleteCompany ...
func (cs *CompanyServiceImpl) DeleteCompany(id int) error {

	err := cs.companyRepo.DeleteCompany(id)

	if err != nil {
		return err
	}

	return nil
}
