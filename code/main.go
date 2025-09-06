package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	ctx := context.Background()
	bot, _ := telego.NewBot(TOKEN)
	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)
	bh, _ := th.NewBotHandler(bot, updates)
	fmt.Print("Бот запущен!")

	defer func() { _ = bh.Stop() }()

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), "Привет! Я бот переводчик написанный на Golang. Введи команду:\n/translate <язык> <текст>\nЧтобы перевести текст."))
		return nil
	}, th.CommandEqual("start"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		args := strings.SplitN(update.Message.Text, " ", 3)
		if len(args) < 3 {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), "Использование: /translate <язык> <текст>"))
			return nil
		}

		language := args[1]
		text := args[2]

		req, _ := http.NewRequest("POST", "https://0.0.0.0:5678/webhook/translate",
			bytes.NewBufferString(fmt.Sprintf(`{"text":"Переведи на %s язык: %s. Строго без контекста"}`, language, text)))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 10 * time.Second}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), string(body)))
		return nil
	}, th.CommandEqual("translate"))
	bh.Start()
}
