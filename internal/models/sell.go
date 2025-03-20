package models

import (
	"fmt"
	"time"
)

type Venda struct {
	ID         int
	Client     *Client
	Session    *Session
	Quantidade int
	DataVenda  time.Time
	ValorTotal float64
}

func NovaVenda(id int, cliente *Client, sessao *Session, quantidade int,
	valorIngresso float64) (*Venda, error) {
	if quantidade <= 0 {
		return nil, fmt.Errorf("quantidade de ingressos deve ser maior que zero")
	}

	if sessao.QuantidadeDisponivel() < quantidade {
		return nil, fmt.Errorf("não há ingressos suficientes disponíveis (%d solicitados, %d disponíveis)",
			quantidade, sessao.QuantidadeDisponivel())
	}

	for i := 0; i < quantidade; i++ {
		if err := sessao.VenderIngresso(); err != nil {
			return nil, err
		}
	}

	return &Venda{
		ID:         id,
		Client:     cliente,
		Session:    sessao,
		Quantidade: quantidade,
		DataVenda:  time.Now(),
		ValorTotal: float64(quantidade) * valorIngresso,
	}, nil
}
