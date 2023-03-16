OpenAI text completion Telegram bot in Golang
======================

Requirements:
- your OpenAPI key in [https://beta.openai.com/account/api-keys](https://beta.openai.com/account/api-keys)
- your Telegram bot token in [https://core.telegram.org/bots/tutorial#obtain-your-bot-token](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)

How to Use
```bash
go run main.go --openai-api-key=xxx --telegram-bot-token=xxx

or 

OPENAI_API_KEY=xxx TELEGRAM_BOT_TOKEN=xxx go run main.go
```