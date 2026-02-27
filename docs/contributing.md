🇧🇷 [Português](#) | 🇺🇸 [English](#english)

---

# Contribuindo

## Configuração de desenvolvimento

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go mod download
```

Build:

```bash
go build -o ezgocommit .
```

Executar sem instalar:

```bash
ANTHROPIC_API_KEY=sk-ant-... go run .
```

## Executando verificações

```bash
go build ./...   # deve compilar sem erros
go vet ./...     # deve passar sem avisos
```

## Estrutura do projeto

```
cmd/            Comandos CLI (Cobra) — sem lógica de negócio
internal/
  config/       carregamento de configuração
  git/          coleta de contexto git (go-git)
  ai/           integração com a API Anthropic
  ui/           TUI Bubbletea
docs/           documentação do projeto
```

Veja [architecture.md](architecture.md) para uma explicação detalhada de cada pacote.

## Adicionando um novo estilo de commit

1. Adicione uma constante em `internal/config/config.go`:
   ```go
   StyleMyStyle = "mystyle"
   ```

2. Documente o estilo em `internal/ai/prompt.go` dentro da constante `systemPrompt` na seção `## Commit styles supported:`.

3. Adicione o estilo à referência de configuração em `docs/configuration.md`.

## Modificando o prompt de IA

O system prompt vive em `internal/ai/prompt.go` como a constante `systemPrompt`. O template do user prompt é `userPromptTemplate` no mesmo arquivo.

Regras para mudanças no prompt:
- O formato de saída deve permanecer JSON estrito correspondendo a `AIResponse` em `internal/ai/types.go`
- Não altere os nomes dos placeholders (`{{GIT_DIFF}}`, etc.) sem atualizar `BuildUserPrompt`
- Teste manualmente com uma chave de API real e uma mudança staged antes de abrir um PR

## Alterando a TUI

A TUI vive inteiramente em `internal/ui/selector.go`. Segue o padrão padrão Bubbletea model/update/view:

- `model` — estado
- `Update()` — trata eventos de teclas, despacha para `updateSelect` ou `updateEdit`
- `View()` — renderiza o estado atual como string usando Lipgloss

Cores são definidas como variáveis `lipgloss.Style` no nível do pacote no topo do arquivo.

## Mensagens de commit

Este projeto usa [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(scope): short description
fix: correct something
docs: update configuration reference
refactor(git): simplify diff collection
```

Tipos: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `style`, `perf`

## Abrindo um pull request

1. Faça um fork e crie um branch: `git checkout -b feat/minha-funcionalidade`
2. Faça suas mudanças
3. Execute `go build ./...` e `go vet ./...`
4. Abra um PR com uma descrição clara da mudança e por que ela foi feita

---

<a id="english"></a>

🇧🇷 [Português](#) | 🇺🇸 [English](#english)

# Contributing

## Development setup

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go mod download
```

Build:

```bash
go build -o ezgocommit .
```

Run without installing:

```bash
ANTHROPIC_API_KEY=sk-ant-... go run .
```

## Running checks

```bash
go build ./...   # must compile with no errors
go vet ./...     # must pass with no warnings
```

## Project layout

```
cmd/            CLI commands (Cobra) — no business logic
internal/
  config/       configuration loading
  git/          git context collection (go-git)
  ai/           Anthropic API integration
  ui/           Bubbletea TUI
docs/           project documentation
```

See [architecture.md](architecture.md) for a detailed explanation of each package.

## Adding a new commit style

1. Add a constant in `internal/config/config.go`:
   ```go
   StyleMyStyle = "mystyle"
   ```

2. Document the style in `internal/ai/prompt.go` inside the `systemPrompt` constant under the `## Commit styles supported:` section.

3. Add the style to the configuration reference in `docs/configuration.md`.

## Modifying the AI prompt

The system prompt lives in `internal/ai/prompt.go` as the `systemPrompt` constant. The user prompt template is `userPromptTemplate` in the same file.

Rules for prompt changes:
- The output format must remain strict JSON matching `AIResponse` in `internal/ai/types.go`
- Do not change placeholder names (`{{GIT_DIFF}}`, etc.) without updating `BuildUserPrompt`
- Test manually with a real API key and a staged change before opening a PR

## Changing the TUI

The TUI lives entirely in `internal/ui/selector.go`. It follows the standard Bubbletea model/update/view pattern:

- `model` — state
- `Update()` — handles key events, dispatches to `updateSelect` or `updateEdit`
- `View()` — renders the current state as a string using Lipgloss

Colors are defined as package-level `lipgloss.Style` variables at the top of the file.

## Commit messages

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(scope): short description
fix: correct something
docs: update configuration reference
refactor(git): simplify diff collection
```

Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `style`, `perf`

## Opening a pull request

1. Fork and create a branch: `git checkout -b feat/my-feature`
2. Make your changes
3. Run `go build ./...` and `go vet ./...`
4. Open a PR with a clear description of the change and why it was made
