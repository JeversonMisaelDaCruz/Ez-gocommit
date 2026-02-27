🇧🇷 [Português](#) | 🇺🇸 [English](#english)

---

# Arquitetura

## Visão geral

O Ez-gocommit é uma ferramenta CLI Go de binário único. Não tem daemon, servidor nem estado persistente além dos arquivos de configuração. Cada invocação executa o pipeline completo e encerra.

## Estrutura do projeto

```
ez-gocommit/
├── main.go                      # Ponto de entrada — chama cmd.Execute()
├── go.mod
├── go.sum
│
├── cmd/
│   ├── root.go                  # Comando raiz Cobra + definição de flags
│   ├── generate.go              # Pipeline principal: coletar → IA → TUI → commit
│   └── version.go               # Subcomando `ezgocommit version`
│
└── internal/
    ├── config/
    │   └── config.go            # Carregar e validar configuração
    │
    ├── git/
    │   └── collector.go         # Coletar contexto git do repositório
    │
    ├── ai/
    │   ├── types.go             # Structs Suggestion e AIResponse
    │   ├── prompt.go            # System prompt + BuildUserPrompt()
    │   └── client.go            # Chamada à API Anthropic + parsing JSON
    │
    └── ui/
        └── selector.go          # TUI interativa com Bubbletea
```

## Fluxo de dados

```
┌─────────────────────────────────────────────────────────────┐
│                        cmd/generate.go                       │
│                                                              │
│  1. config.LoadWithOverrides()                               │
│       ↓                                                      │
│  2. git.Collect()  ──────────────────────────────────────┐  │
│       lê: diff staged, arquivos alterados,               │  │
│           nome do branch, 10 commits recentes, README.md │  │
│       ↓                                                   │  │
│  3. ai.BuildUserPrompt()  ←──────────────────────────────┘  │
│       preenche {{PLACEHOLDERS}} no template do user prompt   │
│       ↓                                                      │
│  4. ai.GenerateSuggestions()                                 │
│       envia system prompt + user prompt para a API Claude    │
│       faz parsing do JSON → []Suggestion                     │
│       ↓                                                      │
│  5. ui.Run()                                                 │
│       TUI Bubbletea: navegar, escolher, editar opcionalmente │
│       retorna (mensagem, corpo, cancelado)                   │
│       ↓                                                      │
│  6. git commit -m "mensagem\n\ncorpo"                        │
└─────────────────────────────────────────────────────────────┘
```

## Responsabilidades dos pacotes

### `internal/config`

Carrega configuração usando Viper. Lê de:
- Arquivo global: `~/.config/ezgocommit/config.toml`
- Arquivo local: `.ezgocommit.toml` (diretório atual)
- Variável de ambiente: `ANTHROPIC_API_KEY`

Flags de CLI aplicadas pelo `cmd` após o carregamento substituem os valores do arquivo.

### `internal/git`

Encapsula `go-git/v5` para extrair tudo necessário para um prompt significativo:

| Função | O que coleta |
|--------|-------------|
| `getBranchName` | Nome curto do branch atual |
| `getStagedDiff` | Diff unificado de HEAD vs index (truncado em `max_diff_lines`) |
| `getRecentCommits` | Últimas 10 linhas de assunto dos commits |
| `getProjectContext` | Primeiras 100 linhas do `README.md` |

Para repositórios sem commits ainda (commit inicial), `getStagedDiff` usa uma lista de arquivos em vez de um patch real.

### `internal/ai`

**`prompt.go`** contém duas constantes string:
- `systemPrompt` — o conjunto completo de instruções enviado como turn de sistema do Claude. Define regras, estilos suportados e o formato de saída JSON estrito.
- `userPromptTemplate` — o template do turn de usuário com tokens `{{PLACEHOLDER}}` substituídos em tempo de execução.

**`client.go`** chama `client.Messages.New()` do `anthropic-sdk-go` oficial. Remove qualquer markdown fence acidental da resposta antes do parsing JSON.

**`types.go`** define `Suggestion` (uma opção) e `AIResponse` (a resposta completa parseada).

### `internal/ui`

Um programa [Bubbletea](https://github.com/charmbracelet/bubbletea) independente com dois modos:

- **`modeSelect`** — navegação com teclas de seta / teclas vim, atalhos numéricos (1-3), `e` para entrar na edição
- **`modeEdit`** — edição inline de texto com movimentação esquerda/direita, `Enter` para confirmar, `Esc` para cancelar

Estilizado com [Lipgloss](https://github.com/charmbracelet/lipgloss). Badges de confiança com código de cores:
- `●● HIGH` → verde
- `●○ MED` → amarelo
- `○○ LOW` → vermelho

### `cmd`

Camada de orquestração fina construída com Cobra. `root.go` registra flags e conecta o comando padrão ao `runGenerate`. `generate.go` possui o pipeline completo em sequência. Nenhuma lógica de negócio vive aqui.

## Dependências

| Pacote | Versão | Papel |
|--------|--------|-------|
| `github.com/spf13/cobra` | v1.10.2 | Estrutura de comandos CLI |
| `github.com/spf13/viper` | v1.21.0 | Carregamento de arquivo de config + variável de ambiente |
| `github.com/go-git/go-git/v5` | v5.17.0 | Operações Git (pure Go) |
| `github.com/anthropics/anthropic-sdk-go` | v1.26.0 | Cliente da API Anthropic |
| `github.com/charmbracelet/bubbletea` | v1.3.10 | Framework de TUI |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | Estilização de TUI |
| `github.com/fatih/color` | v1.18.0 | Saída colorida no terminal |

## Design do prompt de IA

O system prompt instrui o Claude a:
- Analisar o diff profundamente (não apenas nomes de arquivos)
- Usar o nome do branch como dica de intenção
- Espelhar o tom e estilo dos commits recentes
- Entender o domínio do projeto a partir do README
- Produzir exatamente 3 sugestões serializadas em JSON rankeadas por confiança
- Nunca produzir nada fora do objeto JSON

O user prompt encapsula o contexto de runtime em tags XML (`<git_diff>`, `<branch_name>`, etc.) para dar ao Claude limites claros entre cada dado.

## Estratégia de tratamento de erros

- Chave de API ausente → erro claro com instruções de configuração, exit 1
- Sem mudanças staged → erro claro pedindo `git add`, exit 1
- Não é um repositório git → erro do go-git, exit 1
- Erro de API → erro encapsulado com mensagem original, exit 1
- JSON malformado da IA → erro com resposta bruta para debug, exit 1
- Usuário cancela a TUI → imprime "Aborted.", exit 0

---

<a id="english"></a>

🇧🇷 [Português](#) | 🇺🇸 [English](#english)

# Architecture

## Overview

Ez-gocommit is a single-binary Go CLI tool. It has no daemon, no server, and no persistent state beyond config files. Each invocation runs the full pipeline and exits.

## Project structure

```
ez-gocommit/
├── main.go                      # Entry point — calls cmd.Execute()
├── go.mod
├── go.sum
│
├── cmd/
│   ├── root.go                  # Cobra root command + flag definitions
│   ├── generate.go              # Main pipeline: collect → AI → TUI → commit
│   └── version.go               # `ezgocommit version` subcommand
│
└── internal/
    ├── config/
    │   └── config.go            # Load and validate configuration
    │
    ├── git/
    │   └── collector.go         # Collect git context from the repository
    │
    ├── ai/
    │   ├── types.go             # Suggestion and AIResponse structs
    │   ├── prompt.go            # System prompt + BuildUserPrompt()
    │   └── client.go            # Anthropic API call + JSON parsing
    │
    └── ui/
        └── selector.go          # Bubbletea interactive TUI
```

## Data flow

```
┌─────────────────────────────────────────────────────────────┐
│                        cmd/generate.go                       │
│                                                              │
│  1. config.LoadWithOverrides()                               │
│       ↓                                                      │
│  2. git.Collect()  ──────────────────────────────────────┐  │
│       reads: staged diff, changed files,                  │  │
│              branch name, 10 recent commits, README.md    │  │
│       ↓                                                   │  │
│  3. ai.BuildUserPrompt()  ←──────────────────────────────┘  │
│       fills {{PLACEHOLDERS}} in the user prompt template     │
│       ↓                                                      │
│  4. ai.GenerateSuggestions()                                 │
│       sends system prompt + user prompt to Claude API        │
│       parses JSON → []Suggestion                             │
│       ↓                                                      │
│  5. ui.Run()                                                 │
│       Bubbletea TUI: navigate, pick, optionally edit         │
│       returns (message, body, cancelled)                     │
│       ↓                                                      │
│  6. git commit -m "message\n\nbody"                          │
└─────────────────────────────────────────────────────────────┘
```

## Package responsibilities

### `internal/config`

Loads configuration using Viper. Reads from:
- Global file: `~/.config/ezgocommit/config.toml`
- Local file: `.ezgocommit.toml` (current directory)
- Environment variable: `ANTHROPIC_API_KEY`

CLI flags applied by `cmd` after loading override the file values.

### `internal/git`

Wraps `go-git/v5` to extract everything needed for a meaningful prompt:

| Function | What it collects |
|----------|-----------------|
| `getBranchName` | Current branch short name |
| `getStagedDiff` | Unified diff of HEAD vs index (truncated to `max_diff_lines`) |
| `getRecentCommits` | Last 10 commit subject lines |
| `getProjectContext` | First 100 lines of `README.md` |

For repositories with no commits yet (initial commit), `getStagedDiff` falls back to a file list rather than a real patch.

### `internal/ai`

**`prompt.go`** holds two string constants:
- `systemPrompt` — the full instruction set sent as the Claude system turn. Defines rules, supported styles, and the strict JSON output format.
- `userPromptTemplate` — the user turn template with `{{PLACEHOLDER}}` tokens replaced at runtime.

**`client.go`** calls `client.Messages.New()` from the official `anthropic-sdk-go`. It strips any accidental markdown fences from the response before JSON parsing.

**`types.go`** defines `Suggestion` (one option) and `AIResponse` (the full parsed response).

### `internal/ui`

A self-contained [Bubbletea](https://github.com/charmbracelet/bubbletea) program with two modes:

- **`modeSelect`** — arrow key / vim key navigation, number shortcuts (1-3), `e` to enter edit
- **`modeEdit`** — inline text editing with left/right movement, `Enter` to confirm, `Esc` to cancel

Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss). Confidence badges are color-coded:
- `●● HIGH` → green
- `●○ MED` → yellow
- `○○ LOW` → red

### `cmd`

Thin orchestration layer built with Cobra. `root.go` registers flags and wires the default command to `runGenerate`. `generate.go` owns the full pipeline in sequence. No business logic lives here.

## Dependencies

| Package | Version | Role |
|---------|---------|------|
| `github.com/spf13/cobra` | v1.10.2 | CLI command structure |
| `github.com/spf13/viper` | v1.21.0 | Config file + env var loading |
| `github.com/go-git/go-git/v5` | v5.17.0 | Git operations (pure Go) |
| `github.com/anthropics/anthropic-sdk-go` | v1.26.0 | Anthropic API client |
| `github.com/charmbracelet/bubbletea` | v1.3.10 | TUI framework |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | TUI styling |
| `github.com/fatih/color` | v1.18.0 | Colored terminal output |

## AI prompt design

The system prompt instructs Claude to:
- Analyze the diff deeply (not just filenames)
- Use the branch name as an intent hint
- Mirror the tone and style of recent commits
- Understand the project domain from the README
- Produce exactly 3 JSON-serialized suggestions ranked by confidence
- Never output anything outside the JSON object

The user prompt wraps the runtime context in XML-like tags (`<git_diff>`, `<branch_name>`, etc.) to give Claude clear boundaries between each piece of data.

## Error handling strategy

- Missing API key → clear error with setup instructions, exit 1
- No staged changes → clear error prompting `git add`, exit 1
- Not a git repository → error from go-git, exit 1
- API error → wrapped error with original message, exit 1
- Malformed JSON from AI → error with raw response for debugging, exit 1
- User aborts TUI → prints "Aborted.", exit 0
