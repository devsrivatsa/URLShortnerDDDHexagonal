package urlShortner

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("redirect not found")
	ErrRedirectInvalid  = errors.New("redirect is invalid")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

func New(redirectRepo RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepo: redirectRepo,
	}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(err, "service.Redirect.store")
	}
	redirect.Code = shortid.MustGenerate() // generate a short id of 9 characters
	redirect.CreatedAt = time.Now().UTC().Unix()

	return r.redirectRepo.Store(redirect)
}
