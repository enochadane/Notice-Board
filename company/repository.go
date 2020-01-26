package company

import "github.com/amthesonofGod/Notice-Board/entity"

//CompanyRepository ...
type CompanyRepository interface {
	Companies() ([]entity.Company, []error)
	Company(id uint) (*entity.Company, []error)
	CompanyByEmail(email string) (*entity.Company, []error)
	UpdateCompany(company *entity.Company) (*entity.Company, []error)
	DeleteCompany(id uint) (*entity.Company, []error)
	StoreCompany(company *entity.Company) (*entity.Company, []error)
	PhoneExists(phone string) bool
	EmailExists(email string) bool
	// UserRoles(*entity.User) ([]entity.Role, []error)
	
}

//SessionRepositoryCamp ...
type SessionRepositoryCamp interface {
	SessionCamp(sessionID string) (*entity.CompanySession, []error)
	StoreSessionCamp(session *entity.CompanySession) (*entity.CompanySession, []error)
	DeleteSessionCamp(sessionID string) (*entity.CompanySession, []error)
}
