package model

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}
