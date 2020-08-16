package server

import (
	"time"

	"github.com/asdine/storm/v3"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID string `storm:"id"`
	Username string `storm:"unique"`
	CreatedAt time.Time
	HashPassword string
}

type tokenClaims struct {
	jwt.StandardClaims
	username string
}


func (s *Server) performRootAdminCreation(password string) error {
	var userAdmin user
	if err := s.db.One("Username", "admin", &userAdmin); err != nil {
		if err != storm.ErrNotFound {
			logrus.Warn("root admin user not found")
			return errors.WithStack(err)
		}
	}

	if userAdmin.Username == "" { // create new admin
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err  != nil {
			return errors.WithStack(err)
		}
		userAdmin.HashPassword = string(hash)
		userAdmin.Username = "admin"
		userAdmin.ID = uuid.NewV4().String()
		userAdmin.CreatedAt = time.Now()

		logrus.Warn("creating new root admin")
		logrus.Warn("username: admin")
		logrus.Warn("password: " + password)

		if err := s.db.Save(userAdmin); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *Server) login(username, password string) (string, error) {
	var userLogin user
	if err := s.db.One("Username", username, &userLogin); err != nil {
		return "", errors.WithStack(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userLogin.HashPassword), []byte(password)); err != nil {
		return "", errors.WithMessage(err, "invalid username and/or password")
	}

	claims := tokenClaims{
		username: userLogin.Username,
		StandardClaims: jwt.StandardClaims{
			Id: userLogin.ID,
			IssuedAt: time.Now().Unix(),
			Subject: "supersense",
			ExpiresAt: time.Now().Add(1*time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return signedToken, nil
}