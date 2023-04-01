package links

import "context"

type RepositoryInterface interface {
	SaveLink(ctx context.Context, l *Link) error
	GetLatestLinks(ctx context.Context, limit int, page int) ([]Link, error)
	GetLinksFromList(ctx context.Context, listId int, limit int, page int) ([]Link, error)
	GetAllLists(ctx context.Context) ([]List, error)
}
