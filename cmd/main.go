package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GiovannaValadao/circus-tickets/internal/repository"
	service "github.com/GiovannaValadao/circus-tickets/internal/services"
)

func main() {
	repo := repository.NovoRepository()
	svc := service.NovoService(repo, 50.0)

	hoje := time.Now()
	amanha := hoje.AddDate(0, 0, 1)

	_, err := svc.CadastrarSessao(hoje, "15:00")
	if err != nil {
		fmt.Println("Erro ao cadastrar sessão:", err)
	}
	_, err = svc.CadastrarSessao(hoje, "20:00")
	if err != nil {
		fmt.Println("Erro ao cadastrar sessão:", err)
	}

	_, err = svc.CadastrarSessao(amanha, "15:00")
	if err != nil {
		fmt.Println("Erro ao cadastrar sessão:", err)
	}
	_, err = svc.CadastrarSessao(amanha, "20:00")
	if err != nil {
		fmt.Println("Erro ao cadastrar sessão:", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		exibirMenu()
		opcao, _ := reader.ReadString('\n')
		opcao = strings.TrimSpace(opcao)

		switch opcao {
		case "1":
			cadastrarCliente(svc, reader)
		case "2":
			venderIngressos(svc, reader)
		case "3":
			listarSessoes(svc)
		case "4":
			listarClientes(svc)
		case "5":
			gerarRelatorio(svc, reader)
		case "0":
			fmt.Println("Saindo do sistema. Até logo!")
			return
		default:
			fmt.Println("Opção inválida!")
		}

		fmt.Println("\nPressione ENTER para continuar...")
		reader.ReadString('\n')
	}
}

func exibirMenu() {
	limparTela()
	fmt.Println("==== SISTEMA DE GERENCIAMENTO DE INGRESSOS DO CIRCO ====")
	fmt.Println("1. Cadastrar Cliente")
	fmt.Println("2. Vender Ingressos")
	fmt.Println("3. Listar Sessões")
	fmt.Println("4. Listar Clientes")
	fmt.Println("5. Gerar Relatório Diário")
	fmt.Println("0. Sair")
	fmt.Print("Escolha uma opção: ")
}

func limparTela() {
	fmt.Print("\033[H\033[2J")
}

func cadastrarCliente(svc *service.Service, reader *bufio.Reader) {
	limparTela()
	fmt.Println("==== CADASTRO DE CLIENTE ====")

	fmt.Print("Nome: ")
	nome, _ := reader.ReadString('\n')
	nome = strings.TrimSpace(nome)

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Telefone: ")
	telefone, _ := reader.ReadString('\n')
	telefone = strings.TrimSpace(telefone)

	fmt.Print("CPF: ")
	cpf, _ := reader.ReadString('\n')
	cpf = strings.TrimSpace(cpf)

	cliente, err := svc.CadastrarCliente(nome, email, telefone, cpf)
	if err != nil {
		fmt.Println("Erro ao cadastrar cliente:", err)
		return
	}

	fmt.Printf("\nCliente cadastrado com sucesso! ID: %d\n", cliente.ID)
}

func venderIngressos(svc *service.Service, reader *bufio.Reader) {
	limparTela()
	fmt.Println("==== VENDA DE INGRESSOS ====")

	fmt.Println("\nSessões disponíveis:")
	sessoes := svc.ListarSessoes()
	for _, sessao := range sessoes {
		fmt.Printf("ID: %d - %s às %s - Disponíveis: %d\n",
			sessao.ID, sessao.Data.Format("02/01/2006"), sessao.Horario, sessao.QuantidadeDisponivel())
	}

	fmt.Print("\nID da Sessão: ")
	sessaoIDStr, _ := reader.ReadString('\n')
	sessaoIDStr = strings.TrimSpace(sessaoIDStr)
	sessaoID, err := strconv.Atoi(sessaoIDStr)
	if err != nil {
		fmt.Println("ID de sessão inválido!")
		return
	}

	fmt.Print("CPF do Cliente: ")
	cpf, _ := reader.ReadString('\n')
	cpf = strings.TrimSpace(cpf)

	var clienteID int
	clienteExistente, encontrado := svc.BuscarClientePorCPF(cpf)
	if !encontrado {
		fmt.Println("Cliente não encontrado. Vamos cadastrá-lo:")

		fmt.Print("Nome: ")
		nome, _ := reader.ReadString('\n')
		nome = strings.TrimSpace(nome)

		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Print("Telefone: ")
		telefone, _ := reader.ReadString('\n')
		telefone = strings.TrimSpace(telefone)

		cliente, err := svc.CadastrarCliente(nome, email, telefone, cpf)
		if err != nil {
			fmt.Println("Erro ao cadastrar cliente:", err)
			return
		}

		clienteID = cliente.ID
	} else {
		clienteID = clienteExistente.ID
	}

	fmt.Print("Quantidade de ingressos: ")
	quantidadeStr, _ := reader.ReadString('\n')
	quantidadeStr = strings.TrimSpace(quantidadeStr)
	quantidade, err := strconv.Atoi(quantidadeStr)
	if err != nil || quantidade <= 0 {
		fmt.Println("Quantidade inválida!")
		return
	}

	venda, err := svc.RealizarVenda(clienteID, sessaoID, quantidade)
	if err != nil {
		fmt.Println("Erro ao realizar venda:", err)
		return
	}

	fmt.Printf("\nVenda realizada com sucesso! ID: %d - Valor Total: R$ %.2f\n", venda.ID, venda.ValorTotal)
}

func listarSessoes(svc *service.Service) {
	limparTela()
	fmt.Println("==== LISTAGEM DE SESSÕES ====")

	sessoes := svc.ListarSessoes()
	if len(sessoes) == 0 {
		fmt.Println("Nenhuma sessão cadastrada!")
		return
	}

	for _, sessao := range sessoes {
		fmt.Printf("ID: %d - %s às %s - %d/%d ingressos vendidos (%.1f%%)\n",
			sessao.ID,
			sessao.Data.Format("02/01/2006"),
			sessao.Horario,
			sessao.IngressosVendidos,
			sessao.Capacidade,
			float64(sessao.IngressosVendidos)/float64(sessao.Capacidade)*100,
		)
	}
}

func listarClientes(svc *service.Service) {
	limparTela()
	fmt.Println("==== LISTAGEM DE CLIENTES ====")

	clientes := svc.ListarClientes()
	if len(clientes) == 0 {
		fmt.Println("Nenhum cliente cadastrado!")
		return
	}

	for _, cliente := range clientes {
		fmt.Printf("ID: %d - Nome: %s - CPF: %s - Telefone: %s\n",
			cliente.ID, cliente.Name, cliente.CPF, cliente.Phone)
	}
}

func gerarRelatorio(svc *service.Service, reader *bufio.Reader) {
	limparTela()
	fmt.Println("==== RELATÓRIO DIÁRIO ====")

	fmt.Print("Data (DD/MM/AAAA) ou ENTER para hoje: ")
	dataStr, _ := reader.ReadString('\n')
	dataStr = strings.TrimSpace(dataStr)

	var data time.Time
	var err error

	if dataStr == "" {
		data = time.Now()
	} else {
		data, err = time.Parse("02/01/2006", dataStr)
		if err != nil {
			fmt.Println("Formato de data inválido! Use DD/MM/AAAA")
			return
		}
	}

	relatorio := svc.GerarRelatorio(data)
	fmt.Println(relatorio)
}
