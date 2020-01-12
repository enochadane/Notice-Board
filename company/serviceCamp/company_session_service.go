package service

import (
	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
)

// SessionServiceImplCamp implements user.SessionService interface
type SessionServiceImplCamp struct {
	sessionRepo company.SessionRepositoryCamp
}

// NewSessionServiceCamp  returns a new SessionService object
func NewSessionServiceCamp(sessRepository company.SessionRepositoryCamp) company.SessionServiceCamp {
	return &SessionServiceImplCamp{sessionRepo: sessRepository}
}

// SessionCamp returns a given stored session
func (ss *SessionServiceImplCamp) SessionCamp(sessionID string) (*entity.CompanySession, []error) {
	sess, errs := ss.sessionRepo.SessionCamp(sessionID)
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}

// StoreSessionCamp stores a given session
func (ss *SessionServiceImplCamp) StoreSessionCamp(session *entity.CompanySession) (*entity.CompanySession, []error) {
	sess, errs := ss.sessionRepo.StoreSessionCamp(session)
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}

// DeleteSessionCamp deletes a given session
func (ss *SessionServiceImplCamp) DeleteSessionCamp(sessionID string) (*entity.CompanySession, []error) {
	sess, errs := ss.sessionRepo.DeleteSessionCamp(sessionID)
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}
