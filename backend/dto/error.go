package dto

// ErrorResponse は全エンドポイント共通のエラーレスポンス。
// Code はクライアントが機械的に判定するための安定した識別子、Message は人が読むためのメッセージ。
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"error"`
}
