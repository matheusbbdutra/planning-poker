# Planning Poker – Frontend

Aplicação Nuxt 3 que entrega a interface do Planning Poker online. A tela inicial permite criar uma nova sala ou entrar em uma existente sem necessidade de autenticação.

## Tela Inicial

- **Criar sala:** informe um nome para a sessão e seu nome. Ao enviar, uma sala é criada e você entra automaticamente como Scrum Master. O `userId` e o flag `isScrumMaster` são salvos em `sessionStorage`.
- **Entrar em sala:** informe o código da sala e seu nome. O backend retorna seu `userId`, que também é persistido no `sessionStorage`. Você é redirecionado para a URL `/room/{id}` com o WebSocket conectado.
- **Fluxo básico:**
  1. Scrum Master cria a sala e compartilha o `roomId`.
  2. Participantes acessam a tela inicial, escolhem “Entrar em sala” e preenchem o ID.
  3. Todos são levados à tela da sala para votar nas tarefas em tempo real.

### Screenshot da Tela Inicial

> Substitua o caminho abaixo pelo arquivo definitivo quando tiver o print.

![Tela inicial do Planning Poker](./docs/screens/tela-inicial.png)

## Como executar

Instale as dependências e suba o servidor de desenvolvimento em `http://localhost:3000`.

```bash
# instalar dependências
npm install

# iniciar modo desenvolvimento
npm run dev
```

### Build de produção

```bash
# gerar build
npm run build

# pré-visualizar build
npm run preview
```

Para outras opções de gerenciadores de pacote ou detalhes de deploy consulte a [documentação oficial do Nuxt](https://nuxt.com/docs/getting-started/introduction).
