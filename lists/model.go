package lists

type ListAttributes struct {
	FolderId   int    `json:"folder_id"`
	Icon       string `json:"icon"`
	ItemsCount int    `json:"items_count"`
	Name       string `json:"name"`
}

type List struct {
	Type       string         `json:"type"`
	Id         string         `json:"id"`
	Attributes ListAttributes `json:"attributes"`
}

type Lists struct {
	Data []List `json:"data"`
}

type ItemAttributes struct {
	Description  string  `json:"description"`
	IsStarred    bool    `json:"is_starred"`
	ListId       int     `json:"list_id"`
	Name         string  `json:"name"`
	Order        int     `json:"order"`
	Price        float32 `json:"price"`
	Quantity     int     `json:"quantity"`
	QuantityType string  `json:"quantity_type"`
}

type Item struct {
	Type       string         `json:"type"`
	Id         string         `json:"id"`
	Attributes ItemAttributes `json:"attributes"`
}

type Items struct {
	Data []Item `json:"data"`
}
