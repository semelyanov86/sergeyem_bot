package strategies

const NoSettingsError = "В базе данных не найден указанный пользователь. Пожалуйста, авторизируйтесь в чате снова. Для этого введите команду /start 🏁"

const RootModeFailed = "Ошибка при изменении на корневое состояние приложения 😞"

const ContextClearFailed = "Ошибка при очистке контекста 🧐"

const LinkAceTokenError = "Ошибка при авторизации в LinkAce. Пройдите авторизацию заново."

const ContextChangeFailed = "Ошибка при изменении контакста. Что-то не так с нашей БД... 🤬"

const ModeChangeFailed = "Ошибка при переключении режима приложения. Проблема с базой... 🤬"

const (
	Root = iota
	LinksToken
	AskList
	AskListForLinks
)
