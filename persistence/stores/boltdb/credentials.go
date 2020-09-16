package boltdb

import (
	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type credentials struct {
	Username string `storm:"id"`
	Password string
}

// SaveCredential saves a new pair of credentials
func (s *Store) SaveCredential(username, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	c := &credentials{Username: username, Password: string(hashedPass)}
	if err := s.db.Save(c); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ValidateCredential verify if your credentials are ok
func (s *Store) ValidateCredential(username, password string) (bool, error) {
	c := new(credentials)
	if err := s.db.One("Username", username, c); err != nil {
		return false, errors.WithStack(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)); err != nil {
		return false, errors.WithStack(err)
	}

	return true, nil
}

// UpdateCredential helps to update your password
func (s *Store) UpdateCredential(username, password, newPassword string) error {
	isValid, err := s.ValidateCredential(username, password)
	if err != nil {
		return errors.WithStack(err)
	}

	if !isValid {
		return errors.New("invalid credentials")
	}

	newHashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	c := new(credentials)
	if err := s.db.One("Username", username, c); err != nil {
		return errors.WithStack(err)
	}

	c.Password = string(newHashedPass)

	if err = s.db.Update(c); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ForceUpdateCredential forces the inmediatly update of your credential
func (s *Store) ForceUpdateCredential(username, newPassword string) error {
	log.Warnf("Forcing update of your \"%s\" credential.", username)

	newHashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	c := new(credentials)
	if err := s.db.One("Username", username, c); err != nil {
		return errors.WithStack(err)
	}

	c.Password = string(newHashedPass)

	if err = s.db.Update(c); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// UsernameExists check if your username exists
func (s *Store) UsernameExists(username string) (bool, error) {
	c := new(credentials)
	if err := s.db.One("Username", username, c); err != nil {
		if err == storm.ErrNotFound {
			return false, nil
		}
		return false, errors.WithStack(err)
	}

	return true, nil
}
