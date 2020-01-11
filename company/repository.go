package company

import "github.com/amthesonofGod/Notice-Board/entity"

//CompanyRepository ...
type CompanyRepository interface {
	Companies() ([]entity.Company, []error)
	Company(id uint) (*entity.Company, []error)
	UpdateCompany(company *entity.Company) (*entity.Company, []error)
	DeleteCompany(id uint) (*entity.Company, []error)
	StoreCompany(company *entity.Company) (*entity.Company, []error)
	StoreSession(session *entity.CompanySession) (*entity.CompanySession, []error)
	Session(uuid string) (*entity.CompanySession, []error)
	DeleteSession(uuid string) (*entity.CompanySession, []error)
}

//SessionRepositoryCamp ...
type SessionRepositoryCamp interface {
	SessionCamp(sessionID string) (*entity.CompanySession, []error)
	StoreSessionCamp(session *entity.CompanySession) (*entity.CompanySession, []error)
	DeleteSessionCamp(sessionID string) (*entity.CompanySession, []error)
}
