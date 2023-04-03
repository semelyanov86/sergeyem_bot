package links

import (
	"context"
)

type MockRepository struct {
	Token string
}

func (l2 MockRepository) SaveLink(ctx context.Context, l *Link) error {
	return nil
}

func (l2 MockRepository) GetLatestLinks(ctx context.Context, limit int, page int) ([]Link, error) {
	var links = []Link{
		{
			Id:          1,
			URL:         "https://sergeyem.ru",
			Title:       "Sergey Emelyanov",
			Description: "Blog of Software Developer Sergey Emelyanov",
			Icon:        "fa-sergeyem",
			Lists:       nil,
			Tags:        nil,
		},
	}
	if limit > 1 || limit == 0 {
		links = append(links, Link{
			Id:          2,
			URL:         "https://itvolga.com",
			Title:       "Center of Information Technologies",
			Description: "IT Company in Russia for software development",
			Icon:        "fa-itvolga",
			Lists:       nil,
			Tags:        nil,
		})
	}
	return links, nil
}

func (l2 MockRepository) GetLinksFromList(ctx context.Context, listId int, limit int, page int) ([]Link, error) {
	return nil, nil
}

func (l2 MockRepository) GetAllLists(ctx context.Context) ([]List, error) {
	lists := []List{
		{
			Id:          1,
			Name:        "PHP",
			Description: "PHP websites",
			IsPrivate:   false,
		},
		{
			Id:          2,
			Name:        "Golang",
			Description: "Golang websites",
			IsPrivate:   false,
		},
	}
	return lists, nil
}
