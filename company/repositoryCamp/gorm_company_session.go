package repository

import (
	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/jinzhu/gorm"
)

// SessionGormRepoCamp implements user.SessionRepository interface
type SessionGormRepoCamp struct {
	conn *gorm.DB
}

// NewSessionGormRepoCamp  returns a new SessionGormRepo object
func NewSessionGormRepoCamp(db *gorm.DB) company.SessionRepositoryCamp {
	return &SessionGormRepoCamp{conn: db}
}

// SessionCamp returns a given stored session
func (sr *SessionGormRepoCamp) SessionCamp(sessionID string) (*entity.CompanySession, []error) {
	session := entity.CompanySession{}
	errs := sr.conn.Find(&session, "uuid=?", sessionID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &session, errs
}

// StoreSessionCamp stores a given session
func (sr *SessionGormRepoCamp) StoreSessionCamp(session *entity.CompanySession) (*entity.CompanySession, []error) {
	sess := session
	errs := sr.conn.Save(sess).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}

// DeleteSessionCamp deletes a given session
func (sr *SessionGormRepoCamp) DeleteSessionCamp(sessionID string) (*entity.CompanySession, []error) {
	sess, errs := sr.SessionCamp(sessionID)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = sr.conn.Delete(sess, sess.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return sess, errs
}
