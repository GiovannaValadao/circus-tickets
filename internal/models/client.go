package models

type Client struct {
	ID    int
	Name  string
	Email string
	Phone string
	CPF   string
}

func ClientNew(id int, name, email, phone, cpf string) *Client {
	return &Client{
		ID:    id,
		Name:  name,
		Email: email,
		Phone: phone,
		CPF:   cpf,
	}
}
