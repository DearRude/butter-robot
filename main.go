package main

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func main() {
	c := GenConfig()

	ctx := context.Background()

	dispatcher := tg.NewUpdateDispatcher()
	commandHandler := makeHandler()
	opts := telegram.Options{
		Logger:        &c.Logger,
		UpdateHandler: dispatcher,
	}

	client := telegram.NewClient(c.AppId, c.AppHash, opts)

	if err := client.Run(ctx, func(ctx context.Context) error {
		if _, err := client.Auth().Bot(ctx, c.BotToken); err != nil {
			return err
		}

		api := tg.NewClient(client)
		sender := message.NewSender(api)

		dispatcher.OnNewMessage(func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
			m, ok := u.Message.(*tg.Message)
			if !ok || m.Out {
				return nil
			}
			commandHandler.Options = CommandOptions{
				Ctx:      ctx,
				Client:   sender,
				Entities: entities,
				Update:   u,
				Message:  m,
			}
			commandHandler.Run()

			return nil
		})

		if err := telegram.RunUntilCanceled(ctx, client); err != nil {
			c.Logger.Error("Unable for client to run until canceled")
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
