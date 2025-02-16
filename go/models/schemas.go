package models

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []ItemInfo  `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type ItemInfo struct {
	Item     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []RecieveCoinRequest `json:"received"`
	Sent     []SendCoinRequest    `json:"sent"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type RecieveCoinRequest struct {
	Sender string `json:"fromUser"`
	Amount int    `json:"amount"`
}

type SendCoinRequest struct {
	Recipient string `json:"toUser"`
	Amount    int    `json:"amount"`
}
