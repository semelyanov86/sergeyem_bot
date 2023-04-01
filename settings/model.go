package settings

type Setting struct {
	ID             int64
	Username       string
	ChatId         int64
	LinkaceToken   string
	EasylistToken  string
	EasywordsToken string
	Mode           int
	Context        string
}
