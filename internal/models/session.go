package models

import (
	"fmt"
	"time"
)

type Session struct {
	ID                int
	Data              time.Time
	Horario           string
	Capacidade        int
	IngressosVendidos int
}

func NovaSessao(id int, data time.Time, horario string, capacidade int) *Session {
	return &Session{
		ID:                id,
		Data:              data,
		Horario:           horario,
		Capacidade:        capacidade,
		IngressosVendidos: 0,
	}
}

func (s *Session) TemDisponibilidade() bool {
	return s.IngressosVendidos < s.Capacidade
}

func (s *Session) QuantidadeDisponivel() int {
	return s.Capacidade - s.IngressosVendidos
}

func (s *Session) VenderIngresso() error {
	if !s.TemDisponibilidade() {
		return fmt.Errorf("sessão esgotada para o dia %s às %s",
			s.Data.Format("02/01/2006"), s.Horario)
	}

	s.IngressosVendidos++
	return nil
}
