OpenAI text completion Telegram bot in Golang
======================

## Requirements:
- your OpenAPI key in [https://beta.openai.com/account/api-keys](https://beta.openai.com/account/api-keys)
- your Telegram bot token in [https://core.telegram.org/bots/tutorial#obtain-your-bot-token](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)

## How to Use

### Docker
```bash
docker run --restart=always -d \
    -e OPENAI_API_KEY=xxx \
    -e TELEGRAM_BOT_TOKEN=xxx \
    --name chatelegram \
    omegaatt36/chatelegram:latest
```

### Rootless Podman and Systemd
```bash
podman create --restart=always \
    -e OPENAI_API_KEY=xxx \
    -e TELEGRAM_BOT_TOKEN=xxx \
    --name chatelegram \
    omegaatt36/chatelegram:latest
podman generate systemd --new --files --name chatelegram
mkdir -p ~/.config/systemd/user/
cp -Z container-chatelegram.service ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable container-chatelegram.service
systemctl --user start container-chatelegram.service
```

### Source
```bash
go run main.go --openai-api-key=xxx --telegram-bot-token=xxx

or 

OPENAI_API_KEY=xxx TELEGRAM_BOT_TOKEN=xxx go run main.go
```
