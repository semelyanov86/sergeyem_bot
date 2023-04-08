# Sergeyem Bot

Personal telegram bot, written on Golang, which can accept webhook data from third-party services, saving links to LinkAce service and getting words from EasyWords. It can also accept items from EasyList service.

## Installation

Before install, add environment variable - `TELEGRAM_TOKEN`, where you will store your token from bot.
Also copy configuration file to `~/.config/bot.yaml`

```bash
  make build
  make run
```

Or you can simply download precompile script and run it on your server.

## Available commands in chatbot

Telegram bot accepts following commands:
* `/start` - Register user in database setting table.
* `/help` - Showing list of all commands
* `/cancel` - Return to main menu. Changing mode of current user to 0, closing conversation.
* `/links` - Accept number of links param. For example `/links 3` showing 3 links per page. By default 10 links. If you do not enter LinkAce token, bot will ask you for it.
* `/listlinks` - Showing links from specific list. It also accepts per page argument. After entering this command, bot will ask you to select specific list.
* To save link in LinkAce, just enter full URL. Then bot will ask you to select a List.
* `/random` - get list of random words for learning from EasyWords service. You can pass number of words as a param.
* `/saveword` - save new word to easywords service.


## Running Tests

To run tests and code audit processes, run the following command. To enable test database, you need to add env vatiable `BOT_TEST_DB` with dsn db of test database.

```bash
  make audit
```


## Feedback
If you have any feedback, please reach out to us at se@sergeyem.ru

## Support
For support, contact me at telegram: @sergeyem .

## Tech Stack
**Server:** Golang

## Environment Variables
To run this project, you will need to add the following environment variables to your .env file
`TELEGRAM_TOKEN`


## Deployment

To deploy this project run

```bash
  make production/deploy/api
```

## License
[MIT](https://choosealicense.com/licenses/mit/)





