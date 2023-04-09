package lists

import "context"

type RepositoryInterface interface {
	GetLists(ctx context.Context, perPage int) ([]List, error)
	GetItemsFromList(ctx context.Context, listId int, perPage int) ([]Item, error)
}
