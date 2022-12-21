package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionsRepository struct {
	db *gorm.DB
}

func NewSessionsRepository(db *gorm.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	// return nil // TODO: replace this
	return u.db.Create(&session).Error
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	// return nil // TODO: replace this
	return u.db.Where("token = ?", tokenTarget).Delete(&model.Session{}).Error
}

func (u *SessionsRepository) UpdateSessions(session model.Session) error {
	// return nil // TODO: replace this
	return u.db.Table("sessions").Where("username = ?", session.Username).Updates(map[string]interface{}{"token": session.Token, "expiry": session.Expiry}).Error
}

func (u *SessionsRepository) TokenValidity(token string) (model.Session, error) {
	// return model.Session{}, nil // TODO: replace this
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}
	if u.TokenExpired(session) {
		err := u.DeleteSessions(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}

	return session, nil
}

func (u *SessionsRepository) SessionAvailName(name string) (model.Session, error) {
	// return model.Session{}, nil // TODO: replace this
	session := model.Session{}
	err := u.db.Where("username = ?", name).First(&session).Error
	if err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (u *SessionsRepository) SessionAvailToken(token string) (model.Session, error) {
	// return model.Session{}, nil // TODO: replace this
	session := model.Session{}
	err := u.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		fmt.Println(err)
		return session, err
	}
	return session, nil
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
