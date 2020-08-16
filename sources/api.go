package sources

import (
	"context"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
	"github.com/minskylab/supersense"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
)

// API is a simple api fetcher source
type API struct {
	id         string
	sourceName string
	events     *chan supersense.Event
	server *fiber.App
	db *storm.DB
}

type user struct {
	ID string `storm:"id"`
	Username string `storm:"unique"`
	CreatedAt time.Time
	HashPassword string
}

// NewAPI creates and init a new dummy source
func NewAPI(dbPath string) (*API, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	eventsChan := make(chan supersense.Event, 1)

	source := &API{
		id:         uuid.NewV4().String(),
		sourceName: "api",
		events:    &eventsChan,
		server: fiber.New(),
		db: db,
	}
	return source, nil
}

func (a *API) performRootAdminCreation(password string) error {
	var userAdmin user
	if err := a.db.One("Username", "admin", &userAdmin); err != nil {
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

		if err := a.db.Save(userAdmin); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// Run starts the recurrent message issuer
func (a *API) Run(ctx context.Context) error {
	password, _ := gonanoid.Nanoid()
	if err := a.performRootAdminCreation(password); err != nil {
		return errors.WithStack(err)
	}

	a.server.Post("/login", func(c *fiber.Ctx) {

	})

	a.server.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	return nil
}

// Events implements the supersense.Source interface
func (a *API) Events(ctx context.Context) *chan supersense.Event {
	return a.events
}

// Dispose close all streams and flows with the source
func (a *API) Dispose(ctx context.Context) {
	_ = a.db.Close()
}