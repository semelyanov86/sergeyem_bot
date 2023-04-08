package words

type Word struct {
	Id         int    `json:"id,omitempty"`
	Original   string `json:"original,omitempty"`
	Translated string `json:"translated,omitempty"`
	Views      int    `json:"views"`
	Language   string `json:"language,omitempty"`
	Starred    bool   `json:"starred"`
}

type WordSettings struct {
	Paginate        string   `json:"paginate"`
	DefaultLanguage string   `json:"default_language"`
	ShowShared      bool     `json:"show_shared"`
	ShowImported    bool     `json:"show_imported"`
	LanguagesList   []string `json:"languages_list"`
	MainLanguage    string   `json:"main_language"`
}

type WordSettingsResponse struct {
	Data WordSettings `json:"data"`
}
