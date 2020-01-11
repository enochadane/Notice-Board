package service

import (
	"github.com/amthesonofGod/Notice-Board/User"
	"github.com/amthesonofGod/Notice-Board/entity"
)

// SessionServiceImpl implements user.SessionService interface
type SessionServiceImpl struct {
	sessionRepo User.SessionRepository
}

// NewSessionService  returns a new SessionService object
func NewSessionService(sessRepository User.SessionRepository) User.SessionService {
	return &SessionServiceImpl{sessionRepo: sessRepository}
}

// Session returns a given stored session
func (ss *SessionServiceImpl) Session(sessionID string) (*entity.UserSession, []error) {
	sess, errs := ss.sessionRepo.Session(sessionID)
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}

// StoreSession stores a given session
func (ss *SessionServiceImpl) StoreSession(session *entity.UserSession) (*entity.UserSession, []error) {
	sess, errs := ss.sessionRepo.StoreSession(session)
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}

// DeleteSession deletes a given session
func (ss *SessionServiceImpl) DeleteSession(sessionID string) (*entity.UserSession, []error) {
	sess, errs := ss.sessionRepo.DeleteSession(sessionID)
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}
