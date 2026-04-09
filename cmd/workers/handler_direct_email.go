package main

import (
	"fmt"
	"context"

	"github.com/swissymissy/Cardinal/internal/pubsub"
	"github.com/wneessen/go-mail"
)

func (wkrcfg *WorkerConfig) HandlerDirectEmail(event pubsub.DirectEvent) pubsub.AckType {

}

// helper: send email to user
func (wr)