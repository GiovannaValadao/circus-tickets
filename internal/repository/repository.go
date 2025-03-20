package repository

import (
	"fmt"
	"sort"
	"time"

	"github.com/GiovannaValadao/circus-tickets/internal/models"
)

type Repository struct {
	clientes         map[int]*models.Client
	sessoes          map[int]*models.Session
	vendas           map[int]*models.Venda
	proximoIDCliente int
	proximoIDSessao  int
	proximoIDVenda   int
}

func NovoRepository() *Repository {
	return &Repository{
		clientes:         make(map[int]*models.Client),
		sessoes:          make(map[int]*models.Session),
		vendas:           make(map[int]*models.Venda),
		proximoIDCliente: 1,
		proximoIDSessao:  1,
		proximoIDVenda:   1,
	}
}

func (r *Repository) SalvarCliente(nome, email, telefone, cpf string) *models.Client {
	cliente := models.ClientNew(r.proximoIDCliente, nome, email, telefone, cpf)
	r.clientes[cliente.ID] = cliente
	r.proximoIDCliente++
	return cliente
}

func (r *Repository) BuscarClientePorID(id int) (*models.Client, error) {
	cliente, ok := r.clientes[id]
	if !ok {
		return nil, fmt.Errorf("cliente com ID %d não encontrado", id)
	}
	return cliente, nil
}

func (r *Repository) BuscarClientePorCPF(cpf string) (*models.Client, bool) {
	for _, cliente := range r.clientes {
		if cliente.CPF == cpf {
			return cliente, true
		}
	}
	return nil, false
}

func (r *Repository) ListarClientes() []*models.Client {
	clientes := make([]*models.Client, 0, len(r.clientes))
	for _, cliente := range r.clientes {
		clientes = append(clientes, cliente)
	}
	return clientes
}

func (r *Repository) SalvarSessao(data time.Time, horario string, capacidade int) *models.Session {
	sessao := models.NovaSessao(r.proximoIDSessao, data, horario, capacidade)
	r.sessoes[sessao.ID] = sessao
	r.proximoIDSessao++
	return sessao
}

func (r *Repository) BuscarSessaoPorID(id int) (*models.Session, error) {
	sessao, ok := r.sessoes[id]
	if !ok {
		return nil, fmt.Errorf("sessão com ID %d não encontrada", id)
	}
	return sessao, nil
}

func (r *Repository) BuscarSessoesPorData(data time.Time) []*models.Session {
	sessoes := make([]*models.Session, 0)
	for _, sessao := range r.sessoes {
		if sessao.Data.Year() == data.Year() &&
			sessao.Data.Month() == data.Month() &&
			sessao.Data.Day() == data.Day() {
			sessoes = append(sessoes, sessao)
		}
	}

	sort.Slice(sessoes, func(i, j int) bool {
		return sessoes[i].Horario < sessoes[j].Horario
	})

	return sessoes
}

func (r *Repository) ListarSessoes() []*models.Session {
	sessoes := make([]*models.Session, 0, len(r.sessoes))
	for _, sessao := range r.sessoes {
		sessoes = append(sessoes, sessao)
	}

	sort.Slice(sessoes, func(i, j int) bool {
		if !sessoes[i].Data.Equal(sessoes[j].Data) {
			return sessoes[i].Data.Before(sessoes[j].Data)
		}
		return sessoes[i].Horario < sessoes[j].Horario
	})

	return sessoes
}

func (r *Repository) SalvarVenda(cliente *models.Client, sessao *models.Session, quantidade int, valorIngresso float64) (*models.Venda, error) {
	venda, err := models.NovaVenda(r.proximoIDVenda, cliente, sessao, quantidade, valorIngresso)
	if err != nil {
		return nil, err
	}

	r.vendas[venda.ID] = venda
	r.proximoIDVenda++
	return venda, nil
}

func (r *Repository) ListarVendas() []*models.Venda {
	vendas := make([]*models.Venda, 0, len(r.vendas))
	for _, venda := range r.vendas {
		vendas = append(vendas, venda)
	}

	sort.Slice(vendas, func(i, j int) bool {
		return vendas[i].DataVenda.Before(vendas[j].DataVenda)
	})

	return vendas
}
