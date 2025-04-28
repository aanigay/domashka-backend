package sms

type Client struct {
	// some fields
}

func New() *Client {
	return &Client{}
}

func (s *Client) Send(phone, message string) error {
	// some code
	return nil
}
