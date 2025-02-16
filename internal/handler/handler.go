package handler

import (
	AuthHandler "github.com/Azreil-OFD/Avito-test/internal/handler/auth"
	BuyItemHandler "github.com/Azreil-OFD/Avito-test/internal/handler/buyItem"
	InfoHandler "github.com/Azreil-OFD/Avito-test/internal/handler/info"
	SendCoinHandler "github.com/Azreil-OFD/Avito-test/internal/handler/sendCoin"
)

type Handlers struct {
	*AuthHandler.AuthHandler
	*BuyItemHandler.BuyItemHandler
	*InfoHandler.InfoHandler
	*SendCoinHandler.SendCoinHandler
}
