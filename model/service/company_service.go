package service

import (
	"NoticeBoard/entity"
	"NoticeBoard/model"
)

type CompanyServiceImpl struct {
	companyRepo model.CompanyRepository
}

func NewCompanyServiceImpl(CompanyRepos model.CompanyRepository) *CompanyServiceImpl {
	return &CompanyServiceImpl{companyRepo: CompanyRepos}
}

func (cs *CompanyServiceImpl) Companies() ([]entity.Company, error) 

	companies, err := cs.companyRepo.Companies()

	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (cs *CompanyServiceImpl) StoreCompany(company entity.Company) error {

	err := cs.companyRepo.StoreCompany(company)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CompanyServiceImpl) Company(id int) (entity.Company, error) {

	c, err := cs.companyRepo.Company(id)

	if err != nil {
		return c, err
	}

	return c, nil
}

func (cs *CompanyServiceImpl) UpdateCompany(company entity.Company) error {

	err := cs.companyRepo.UpdateCompany(company)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CompanyServiceImpl) DeleteCompany(id int) error {

	err := cs.companyRepo.DeleteCompany(id)

	if err != nil {
		return err
	}

	return nil
}