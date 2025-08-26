package main

import (
	"context"
	"log"

	"github.com/denkhaus/agents/di"
	"github.com/denkhaus/agents/provider/chat"
	"github.com/denkhaus/agents/system/multi"
	"github.com/google/uuid"
	"github.com/samber/do"
)

func startup(ctx context.Context) error {

	injector := di.NewContainer()
	multiAgentSystem := do.MustInvoke[multi.MultiAgentSystem](injector)

	chat, err := multiAgentSystem.CreateProjectManagerChat(ctx, injector,
		chat.WithAppName("denkhaus-system-chat"),
		chat.WithSessionID(uuid.New()),
		chat.WithUserID(uuid.New()),
	)

	if err != nil {
		return err
	}

	return multiAgentSystem.EnterChat(ctx, chat)
}

func main() {
	if err := startup(context.Background()); err != nil {
		log.Fatal(err)
	}
}
