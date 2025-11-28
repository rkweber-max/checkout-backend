# 🐳 Docker Setup

Este guia mostra como executar a aplicação usando Docker.

## Pré-requisitos

- Docker
- Docker Compose

## Como Executar

### 1. Parar servidores locais (se estiverem rodando)

```bash
# Parar qualquer instância do Go rodando localmente
# Ctrl+C nos terminais ou:
pkill -f "go run cmd/app/main.go"
```

### 2. Iniciar com Docker Compose

```bash
docker-compose up --build
```

Isso irá:
- Criar a imagem da aplicação Go
- Iniciar o container PostgreSQL
- Criar o banco de dados `productsdb`
- Executar o script `init.sql` para criar a tabela `products`
- Iniciar a aplicação na porta 8080

### 3. Testar a aplicação

```bash
# Criar produto
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Docker Product","description":"Created in Docker","price":99.99}'

# Listar produtos
curl http://localhost:8080/api/products
```

## Comandos Úteis

```bash
# Iniciar em background
docker-compose up -d

# Ver logs
docker-compose logs -f app

# Parar containers
docker-compose down

# Parar e remover volumes (limpa o banco)
docker-compose down -v

# Reconstruir imagens
docker-compose build --no-cache
```

## Estrutura

- **Dockerfile**: Multi-stage build para otimizar o tamanho da imagem
- **docker-compose.yml**: Orquestra app + PostgreSQL
- **init.sql**: Script de inicialização do banco
- **.dockerignore**: Exclui arquivos desnecessários da imagem

## Portas

- **8080**: Aplicação Go
- **5432**: PostgreSQL

## Variáveis de Ambiente

Configuradas no `docker-compose.yml`:
- `APP_PORT=8080`
- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=maxter`
- `DB_PASSWORD=admin`
- `DB_NAME=productsdb`
- `DB_SSLMODE=disable`
