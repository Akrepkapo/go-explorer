package models

import (
	"context"
	"time"

	"github.com/IBAX-io/go-explorer/conf"
)

var centrifugoTimeout = time.Second * 5

const (
	ChannelTopData         = "topdata"
	ChannelBlockAndTxsList = "blocktransactionlist"
)

func WriteChannelByte(channel string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), centrifugoTimeout)
	defer cancel()
