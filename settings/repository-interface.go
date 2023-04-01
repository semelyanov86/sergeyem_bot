package settings

import "errors"

var ErrNoRecord = errors.New("settings: no matching record found")

type RepositoryInterface interface {
	Get(user string) (*Setting, error)
	Insert(setting *Setting) error
	Update(setting *Setting) error
	UpdateMode(userName string, mode int) error
	UpdateContext(userName string, context string) error
}
