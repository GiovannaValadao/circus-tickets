package service

import (
	"fmt"
	"time"

	"github.com/GiovannaValadao/circus-tickets/internal/models"
	"github.com/GiovannaValadao/circus-tickets/internal/repository"
)

type Service struct {
	repo          *repository.Repository
	precoIngresso float64
}

func NovoService(repo *repository.Repository, precoIngresso float64) *Service {
	return &Service{
		repo:          repo,
		precoIngresso: precoIngresso,
	}
}

func (s *Service) CadastrarCliente(nome, email, telefone, cpf string) (*models.Client, error) {
	if clienteExistente, encontrado := s.repo.BuscarClientePorCPF(cpf); encontrado {
		return clienteExistente, fmt.Errorf("cliente com CPF %s já cadastrado (ID: %d)",
			cpf, clienteExistente.ID)
	}

	return s.repo.SalvarCliente(nome, email, telefone, cpf), nil
}

func (s *Service) BuscarClientePorCPF(cpf string) (*models.Client, bool) {
	return s.repo.BuscarClientePorCPF(cpf)
}

func (s *Service) CadastrarSessao(data time.Time, horario string) (*models.Session, error) {
	sessoesDoDia := s.repo.BuscarSessoesPorData(data)
	if len(sessoesDoDia) >= 2 {
		return nil, fmt.Errorf("já existem duas sessões agendadas para o dia %s",
			data.Format("02/01/2006"))
	}

	for _, sessao := range sessoesDoDia {
		if sessao.Horario == horario {
			return nil, fmt.Errorf("já existe uma sessão agendada para %s às %s",
				data.Format("02/01/2006"), horario)
		}
	}

	return s.repo.SalvarSessao(data, horario, 200), nil
}

func (s *Service) RealizarVenda(clienteID, sessaoID, quantidade int) (*models.Venda, error) {
	// Busca cliente
	cliente, err := s.repo.BuscarClientePorID(clienteID)
	if err != nil {
		return nil, err
	}

	sessao, err := s.repo.BuscarSessaoPorID(sessaoID)
	if err != nil {
		return nil, err
	}

	if sessao.QuantidadeDisponivel() < quantidade {
		return nil, fmt.Errorf("não há ingressos suficientes disponíveis (%d solicitados, %d disponíveis)",
			quantidade, sessao.QuantidadeDisponivel())
	}

	return s.repo.SalvarVenda(cliente, sessao, quantidade, s.precoIngresso)
}

func (s *Service) GerarRelatorio(data time.Time) string {
	sessoes := s.repo.BuscarSessoesPorData(data)
	if len(sessoes) == 0 {
		return fmt.Sprintf("Não há sessões agendadas para o dia %s", data.Format("02/01/2006"))
	}

	relatorio := fmt.Sprintf("RELATÓRIO DE OCUPAÇÃO - %s\n\n", data.Format("02/01/2006"))

	totalIngressos := 0
	totalCapacidade := 0

	for _, sessao := range sessoes {
		relatorio += fmt.Sprintf("Sessão das %s: %d/%d ocupação (%.1f%%)\n",
			sessao.Horario,
			sessao.IngressosVendidos,
			sessao.Capacidade,
			float64(sessao.IngressosVendidos)/float64(sessao.Capacidade)*100,
		)

		totalIngressos += sessao.IngressosVendidos
		totalCapacidade += sessao.Capacidade
	}

	relatorio += fmt.Sprintf("\nOcupação Total do Dia: %d/%d (%.1f%%)\n",
		totalIngressos,
		totalCapacidade,
		float64(totalIngressos)/float64(totalCapacidade)*100,
	)

	return relatorio
}

func (s *Service) ListarSessoes() []*models.Session {
	return s.repo.ListarSessoes()
}

func (s *Service) ListarClientes() []*models.Client {
	return s.repo.ListarClientes()
}

func (s *Service) ListarVendas() []*models.Venda {
	return s.repo.ListarVendas()
}
