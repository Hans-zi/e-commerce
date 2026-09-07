package dto

type ChatMsg struct {
	Msg string `json:"msg"`
}

type ChatRes struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
