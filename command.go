package main

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type CommandHandler struct {
	Prefix   string
	Options  CommandOptions
	Commands []BotCommmand
}

type CommandOptions struct {
	Ctx      context.Context
	Client   *message.Sender
	Entities tg.Entities
	Update   *tg.UpdateNewMessage
	Message  *tg.Message
}

type BotCommmand struct {
	Command     string
	Description string
	Func        func(CommandOptions) error
}

func (h *CommandHandler) AddCommand(command string, description string, f func(o CommandOptions) error) {
	botCommand := BotCommmand{
		Command:     command,
		Description: description,
		Func:        f,
	}
	h.Commands = append(h.Commands, botCommand)
}

func (h *CommandHandler) Run() error {
	headCommand := strings.Fields(h.Options.Message.Message)[0]

	for _, command := range h.Commands {
		if headCommand != (h.Prefix + command.Command) {
			continue
		}
		return command.Func(h.Options)
	}

	return nil
}

func makeHandler() CommandHandler {
	handler := CommandHandler{
		Prefix: "/",
	}

	handler.AddCommand("ping", "ping bot availbility", func(o CommandOptions) error {
		_, err := o.Client.Reply(o.Entities, o.Update).Text(o.Ctx, "pong")
		return err
	})

	handler.AddCommand("start", "bot start message", func(o CommandOptions) error {
		text := "I'm ready to pass the butter..."
		_, err := o.Client.Reply(o.Entities, o.Update).Text(o.Ctx, text)
		return err
	})

	handler.AddCommand("uuid", "generate v4 UUID", func(o CommandOptions) error {
		id_string := uuid.New().String()
		_, err := o.Client.Reply(o.Entities, o.Update).Text(o.Ctx, id_string)
		return err
	})

	return handler
}
