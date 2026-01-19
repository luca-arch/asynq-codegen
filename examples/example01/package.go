package example01

//go:generate go run ../..

// asynq:task
type SendEmail struct {
	Body string `json:"body"`
	From string `json:"from"`
	To   string `json:"to"`
}

// asynq:task    send_sms_message
type SendSMS struct {
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}
