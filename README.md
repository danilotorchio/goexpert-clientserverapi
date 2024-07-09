# GoExpert - ClientServerAPI

Projeto relativo ao desafio Client-Server-API durante a **Pós-Graduação em Desenvolvimento Avançado em Go da Faculdade Full Cycle de Tecnologia (FCTECH)**, turma de 2024.


## ⚡️ Execução do projecto

Execute os comandos abaixo para clonar o repositório, baixar as dependências e rodar os dois sistemas:

### 🖥️ Clone do projeto:

```bash
git clone https://github.com/danilotorchio/goexpert-clientserverapi.git
```

### 🖥️ Servidor:

```bash
cd goexpert-clientserverapi/server
go mod tidy
go run server.go
```

### 🖥️ Cliente:

```bash
cd ../client
go run client.go
```
## ⚙️ Sobre a execução

O banco de dados (SQLite) será inicializado automaticamente na pasta database do servidor.

Todos os log's relativos aos erros do sistema são apresentados no console. Caso nenhum erro aconteça, um arquivo *cotacao.txt* será criado com a resposta do servidor (API) na mesma pasta do cliente.

## 🤖 Bonus

Também foi criado um novo endpoint, além do solicitado no desafio, para obter o histórico de cotações realizadas.

- Endpoint: /cotacao
- Method: GET
- Query: limit - integer - Máx: 100
- Response: Array\<cotacao\>
