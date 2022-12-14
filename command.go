package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

type CommandHandler struct {
	Prefix   string
	Logger   *zap.Logger
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
		h.Logger.Info("Bot command is called.", zap.String("command", command.Command))
		return command.Func(h.Options)
	}

	return nil
}

func makeHandler(logger *zap.Logger) CommandHandler {
	handler := CommandHandler{
		Prefix: "/",
		Logger: logger,
	}

	handler.AddCommand("ping", "ping bot availbility", func(o CommandOptions) error {
		tm := time.Now()
		_, err := o.Client.Reply(o.Entities, o.Update).StyledText(o.Ctx,
			styling.Italic(fmt.Sprintf("pong %dμs", time.Since(tm).Microseconds())))
		return err
	})

	handler.AddCommand("start", "bot start message", func(o CommandOptions) error {
		text := "I'm ready to pass the butter..."
		_, err := o.Client.Reply(o.Entities, o.Update).Text(o.Ctx, text)
		return err
	})

	handler.AddCommand("uuid", "generate v4 UUID", func(o CommandOptions) error {
		id_string := uuid.New().String()
		_, err := o.Client.Reply(o.Entities, o.Update).StyledText(o.Ctx, styling.Code(id_string))
		return err
	})

	return handler
}
