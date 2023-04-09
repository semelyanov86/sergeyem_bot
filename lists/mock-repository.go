package lists

import (
	"context"
)

type MockRepository struct {
	Token string
}

func (l MockRepository) GetLists(ctx context.Context, perPage int) ([]List, error) {
	return []List{
		{
			Type: "lists",
			Id:   "1",
			Attributes: ListAttributes{
				FolderId:   1,
				Icon:       "mdi-list",
				ItemsCount: 3,
				Name:       "Soupe Schi",
			},
		},
		{
			Type: "lists",
			Id:   "2",
			Attributes: ListAttributes{
				FolderId:   1,
				Icon:       "mdi-list",
				ItemsCount: 5,
				Name:       "Borsch",
			},
		},
	}, nil
}

func (l MockRepository) GetItemsFromList(ctx context.Context, listId int, perPage int) ([]Item, error) {
	return []Item{
		{
			Type: "items",
			Id:   "4",
			Attributes: ItemAttributes{
				Description:  "Some kapusta",
				IsStarred:    false,
				ListId:       1,
				Name:         "Kapusta",
				Order:        1,
				Price:        10,
				Quantity:     1,
				QuantityType: "st",
			},
		}, {
			Type: "items",
			Id:   "5",
			Attributes: ItemAttributes{
				Description:  "Svekla krasn",
				IsStarred:    false,
				ListId:       2,
				Name:         "Svekla",
				Order:        2,
				Price:        20,
				Quantity:     2,
				QuantityType: "st",
			},
		}, {
			Type: "items",
			Id:   "6",
			Attributes: ItemAttributes{
				Description:  "Kartoshka svezhaya",
				IsStarred:    false,
				ListId:       1,
				Name:         "Kartofel",
				Order:        4,
				Price:        10,
				Quantity:     1,
				QuantityType: "kg",
			},
		},
	}, nil
}
