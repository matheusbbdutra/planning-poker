# Planning Poker

Aplicação completa para facilitar estimativas de times ágeis usando a técnica de Planning Poker. O projeto é composto por uma API escrita em Go, frontend em Nuxt 3 e Redis para comunicação em tempo real. Há ambientes prontos para desenvolvimento e produção utilizando Docker Compose.

## Visão geral da arquitetura

- **API (Go 1.25)** – expõe endpoints REST para criação/entrada em salas (`/v1/room`) e mantém um hub WebSocket para broadcast dos eventos de votação.
- **Frontend (Nuxt 3 + Tailwind)** – interface responsiva para criar ou entrar em salas, gerenciar tarefas, cadastrar baralhos e acompanhar as médias de votos.
- **Redis** – usado como canal de publicação/assinatura e armazenamento temporário das salas.
- **Traefik (produção)** – roteador reverso responsável por TLS e roteamento entre API e frontend.

## Funcionalidades principais

- Criação de sala com identificação única e definição automática do Scrum Master.
- Convite de participantes por código e persistência do `userId` no `sessionStorage`.
- Votação em tempo real via WebSocket, com revelação de votos e cálculo de média.
- Cadastro e gestão de tarefas dentro da sala, além de personalização das cartas de votação.
- Deploy orientado a containers, com imagens separadas para backend, frontend e Redis.

## Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) e Docker Compose
- Opcional para execução manual:
  - Go 1.25+
  - Node.js 20+ ou 18 LTS
  - Redis 7+

## Executando com Docker Compose

Com todos os serviços orquestrados (modo recomendado para desenvolvimento):

```bash
# subir ambiente de desenvolvimento
make dev-up

# ou diretamente
# docker compose -f docker-compose.dev.yml up --build

# encerrar os serviços
make dev-down
```

Os serviços disponíveis neste modo são:

- **API** acessível em `http://localhost:8081/v1`
- **Frontend** em `http://localhost:3001`
- **Redis** exposto na porta `6379`

Para ambiente de produção (build otimizado, sem bind mounts):

```bash
make prod-up
# docker compose -f docker-compose.prod.yml up -d --build
```

O arquivo `docker-compose.prod.yml` pressupõe que exista uma rede externa do Traefik e os resolvers configurados (`cloudflare`). Ajuste conforme a sua infraestrutura.

## Executando os serviços manualmente

### 1. Redis

```bash
docker run --name redis-planning -p 6379:6379 -d redis:alpine
```

### 2. API em Go

```bash
cd api
REDIS_ADDR=localhost:6379 go run ./cmd/app
```

A API sobe por padrão na porta `8080`.

### 3. Frontend Nuxt

```bash
cd web
npm install
API_URL=http://localhost:8081/v1 npm run dev
```

O arquivo `web/.env` já traz o valor padrão de `API_URL` apontando para a porta exposta pelo compose de desenvolvimento.

## Estrutura de pastas

```
.
├── api/                  # código fonte da API em Go
├── web/                  # aplicação Nuxt 3
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── Makefile              # atalhos para os comandos Docker Compose
└── README.md             # este arquivo
```

## Fluxo de desenvolvimento

1. O facilitador cria uma sala na interface e recebe o `roomId`.
2. Participantes entram com o código e um nome, sendo adicionados via endpoint `/room/{roomId}/join`.
3. O hub WebSocket (`/v1/room/ws/{roomId}/{userId}`) mantém todos sincronizados durante as rodadas de votação.
4. Após revelar os votos, a média é calculada no backend e devolvida para os clientes conectados.

## Links úteis

- Frontend ao vivo: https://planning-poker.mdutra.dev
- Repositório do projeto: https://github.com/matheusbbdutra/planning-poker
- Documentação Nuxt: https://nuxt.com/docs
- Documentação Go: https://go.dev/doc/

## Contribuindo

Fique à vontade para abrir *issues* ou enviar *pull requests*. Sugestões de novas cartas, melhorias de UX ou monitoramento são bem-vindas.

---

© 2025 Matheus Dutra. Atualize este aviso conforme necessário.
