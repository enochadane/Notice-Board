package company

import "github.com/amthesonofGod/Notice-Board/entity"

// CompanyService specifies compny services

//CompanyService ...

type CompanyService interface {
	Companies() ([]entity.Company, []error)
	Company(id uint) (*entity.Company, []error)
	CompanyByEmail(email string) (*entity.Company, []error)
	UpdateCompany(company *entity.Company) (*entity.Company, []error)
	DeleteCompany(id uint) (*entity.Company, []error)
	StoreCompany(company *entity.Company) (*entity.Company, []error)
	PhoneExists(phone string) bool
	EmailExists(email string) bool
	CompanyRoles(company *entity.Company) ([]entity.Role, []error)
}

// RoleService speifies application user role related database operations
type RoleService interface {
	Roles() ([]entity.Role, []error)
	Role(id uint) (*entity.Role, []error)
	RoleByName(name string) (*entity.Role, []error)
	UpdateRole(role *entity.Role) (*entity.Role, []error)
	DeleteRole(id uint) (*entity.Role, []error)
	StoreRole(role *entity.Role) (*entity.Role, []error)
}

//SessionServiceCamp ...
type SessionServiceCamp interface {
	SessionCamp(sessionID string) (*entity.CompanySession, []error)
	StoreSessionCamp(session *entity.CompanySession) (*entity.CompanySession, []error)
	DeleteSessionCamp(sessionID string) (*entity.CompanySession, []error)
}
