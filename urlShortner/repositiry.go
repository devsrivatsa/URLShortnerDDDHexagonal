package urlShortner

type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}