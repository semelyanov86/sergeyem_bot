package links

import (
	"encoding/json"
	"strings"
)

type Link struct {
	Id          int    `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Lists       Lists  `json:"lists,omitempty"`
	Tags        []Tag  `json:"tags,omitempty"`
}

type List struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private,omitempty"`
}

type Lists []List

type Tag struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (lists *Lists) MarshalJSON() ([]byte, error) {
	var names []string
	for _, list := range *lists {
		names = append(names, list.Name)
	}
	return json.Marshal(strings.Join(names, ","))
}

func (lists *Lists) UnmarshalJSON(data []byte) error {
	var listModels []List
	if err := json.Unmarshal(data, &listModels); err != nil {
		return err
	}
	for _, list := range listModels {
		*lists = append(*lists, list)
	}
	return nil
}
