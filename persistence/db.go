package persistence

import (
	"math/rand"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Persistence struct {
	db *storm.DB
}

func New(dbPath string) (*Persistence, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Persistence{
		db: db,
	}, nil
}

// User represents any user that interact with supersense
type User struct {
	ID string `storm:"id"`
	Username string `storm:"unique"`
	CreatedAt time.Time
	HashPassword string
}

func (db *Persistence) performRootAdminCreation(password string) error {
	var userAdmin User
	if err := db.db.One("Username", "admin", &userAdmin); err != nil {
		if err != storm.ErrNotFound {
			logrus.Warn("root admin User not found")
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

		if err := db.db.Save(&userAdmin); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// ValidateIfAdminExist perform two steps, first validate if root admin exists and if not
// creates a new admin root with a alpha-numeric password with 10 characters (simple and insecure)
func (db *Persistence) ValidateIfAdminExist() error {
	rand.Seed(time.Now().Unix())
	passwordLength := 10
	alphabet := "abcdefghijklmnopqrst1234567890"
	password := ""
	for len(password) < passwordLength{
		i := rand.Intn(len(alphabet))
		password += string(alphabet[i])
	}

	return db.performRootAdminCreation(password)
}

func (db *Persistence) login(username, password string) (*User, error) {
	userLogin := new(User)
	if err := db.db.One("Username", username, userLogin); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userLogin.HashPassword), []byte(password)); err != nil {
		return nil, errors.WithMessage(err, "invalid username and/or password")
	}
	userLogin.HashPassword = ""
	return userLogin, nil
}

// LoginWithUserPassword perform a simple query to persistence and compare its saved hash
func (db *Persistence) LoginWithUserPassword(username, password string) (*User, error) {
	return db.login(username, password)
}

// GetUserByID find one user with ID passed as a param
func (db *Persistence) GetUserByID(id string) (*User, error) {
	user := new(User)
	if err := db.db.One("ID", id, user); err != nil {
		return nil, errors.WithStack(err)
	}
	user.HashPassword = ""
	return user, nil
}