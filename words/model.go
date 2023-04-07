package words

type Word struct {
	Id         int    `json:"id,omitempty"`
	Original   string `json:"original,omitempty"`
	Translated string `json:"translated,omitempty"`
	Views      int    `json:"views,omitempty"`
	Language   string `json:"language,omitempty"`
	Starred    bool   `json:"starred,omitempty"`
}
