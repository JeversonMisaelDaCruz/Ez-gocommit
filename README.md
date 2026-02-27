🇧🇷 [Português](#) | 🇺🇸 [English](#english)

---

# Ez-gocommit

Uma ferramenta CLI escrita em Go que gera mensagens de commit Git semânticas usando Claude (Anthropic) ou Gemini (Google). Ela analisa seu diff staged, nome do branch, histórico de commits recentes e o README do projeto para produzir 3 sugestões rankeadas — exibidas em uma TUI interativa no terminal onde você pode escolher, editar ou cancelar.

O provedor é detectado automaticamente pela chave de API: chaves `AIzaSy*` usam Gemini; qualquer outra usa Claude.

## Como funciona

```
git add .  →  ezgocommit  →  [TUI com 3 sugestões rankeadas]  →  git commit
```

1. Lê seu diff staged, arquivos alterados, nome do branch, commits recentes e `README.md`
2. Envia esse contexto para o Claude via API Anthropic
3. Retorna 3 mensagens de commit rankeadas (alta / média / baixa confiança)
4. Permite navegar, escolher ou editar inline — e então commita automaticamente

## Instalação

**A partir do código-fonte:**

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go build -o ezgocommit .
sudo mv ezgocommit /usr/local/bin/
```

**Com `go install`:**

```bash
go install github.com/jeversonmisael/ez-gocommit@latest
```

## Requisitos

- Go 1.22+
- Uma chave de API: [Anthropic](https://console.anthropic.com/) ou [Google AI Studio](https://aistudio.google.com/)
- Um repositório Git com mudanças staged

## Configuração

A única configuração obrigatória é sua chave de API. O provedor é detectado automaticamente pelo prefixo da chave.

**Com Claude (Anthropic):**

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

**Com Gemini (Google):**

```bash
export ANTHROPIC_API_KEY=AIzaSy...
```

Adicione ao seu `~/.zshrc` ou `~/.bashrc` para persistir. Também é possível usar um arquivo de configuração:

```toml
# ~/.config/ezgocommit/config.toml
api_key = "sk-ant-..."   # ou AIzaSy... para Gemini
```

Veja [docs/configuration.md](docs/configuration.md) para todas as opções disponíveis.

## Uso

```bash
git add .
ezgocommit
```

```
⠸ Analyzing your changes with Claude...

╭──────────────────────────────────────────────────────────────────────╮
│  Ez-gocommit — Select a commit message                               │
│                                                                      │
│  ▶ [1] ●● HIGH   feat(auth): add JWT refresh token rotation          │
│    [2] ●○ MED    feat(auth): implement token refresh endpoint        │
│    [3] ○○ LOW    chore(auth): update token handling logic            │
│                                                                      │
│  💬 Branch name and diff clearly indicate authentication token logic │
│                                                                      │
│  ↑↓/jk navigate • 1-3 jump • Enter confirm • e edit • q abort       │
╰──────────────────────────────────────────────────────────────────────╯

✔ Committed: feat(auth): add JWT refresh token rotation
```

### Controles da TUI

| Tecla | Ação |
|-------|------|
| `↑` / `↓` ou `j` / `k` | Navegar entre sugestões |
| `1` / `2` / `3` | Ir diretamente para aquela sugestão |
| `Enter` | Confirmar e commitar |
| `e` | Editar a mensagem selecionada inline |
| `q` / `Esc` / `Ctrl+C` | Cancelar |

**No modo de edição:**

| Tecla | Ação |
|-------|------|
| `Enter` | Confirmar mensagem editada |
| `Esc` | Cancelar edição, voltar à seleção |
| `←` / `→` | Mover cursor |
| `Ctrl+A` / `Home` | Ir para o início |
| `Ctrl+E` / `End` | Ir para o fim |
| `Backspace` | Deletar caractere |

### Flags

```bash
ezgocommit --style gitmoji       # usar gitmoji em vez de conventional commits
ezgocommit --style free          # sem restrições de formato
ezgocommit --model claude-opus-4-6  # usar um modelo Claude diferente
```

## Estilos de commit

| Estilo | Exemplo |
|--------|---------|
| `conventional` (padrão) | `feat(auth): add login with OAuth` |
| `gitmoji` | `✨ add login with OAuth` |
| `free` | `Add OAuth login support` |
| `custom` | Definido por você em `custom_format` |

## Documentação

- [Primeiros Passos](docs/getting-started.md)
- [Referência de Configuração](docs/configuration.md)
- [Arquitetura](docs/architecture.md)
- [Contribuindo](docs/contributing.md)

## Licença

MIT

---

<a id="english"></a>

🇧🇷 [Português](#) | 🇺🇸 [English](#english)

# Ez-gocommit

A CLI tool written in Go that generates semantic Git commit messages using Claude (Anthropic) or Gemini (Google). It analyzes your staged diff, branch name, recent commit history, and project README to produce 3 ranked suggestions — displayed in an interactive terminal UI where you can pick, edit, or abort.

The provider is detected automatically from your API key: `AIzaSy*` keys use Gemini; anything else uses Claude.

## How it works

```
git add .  →  ezgocommit  →  [TUI with 3 ranked suggestions]  →  git commit
```

1. Reads your staged diff, changed files, branch name, recent commits, and `README.md`
2. Sends that context to Claude via the Anthropic API
3. Returns 3 ranked commit messages (high / medium / low confidence)
4. Lets you navigate, pick, or edit one inline — then commits automatically

## Install

**From source:**

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go build -o ezgocommit .
sudo mv ezgocommit /usr/local/bin/
```

**With `go install`:**

```bash
go install github.com/jeversonmisael/ez-gocommit@latest
```

## Requirements

- Go 1.22+
- An API key: [Anthropic](https://console.anthropic.com/) or [Google AI Studio](https://aistudio.google.com/)
- A Git repository with staged changes

## Setup

The only required configuration is your API key. The provider is detected automatically from the key prefix.

**With Claude (Anthropic):**

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

**With Gemini (Google):**

```bash
export ANTHROPIC_API_KEY=AIzaSy...
```

Add it to your `~/.zshrc` or `~/.bashrc` to persist it. You can also use a config file:

```toml
# ~/.config/ezgocommit/config.toml
api_key = "sk-ant-..."   # or AIzaSy... for Gemini
```

See [docs/configuration.md](docs/configuration.md) for all available options.

## Usage

```bash
git add .
ezgocommit
```

```
⠸ Analyzing your changes with Claude...

╭──────────────────────────────────────────────────────────────────────╮
│  Ez-gocommit — Select a commit message                               │
│                                                                      │
│  ▶ [1] ●● HIGH   feat(auth): add JWT refresh token rotation          │
│    [2] ●○ MED    feat(auth): implement token refresh endpoint        │
│    [3] ○○ LOW    chore(auth): update token handling logic            │
│                                                                      │
│  💬 Branch name and diff clearly indicate authentication token logic │
│                                                                      │
│  ↑↓/jk navigate • 1-3 jump • Enter confirm • e edit • q abort       │
╰──────────────────────────────────────────────────────────────────────╯

✔ Committed: feat(auth): add JWT refresh token rotation
```

### TUI controls

| Key | Action |
|-----|--------|
| `↑` / `↓` or `j` / `k` | Navigate between suggestions |
| `1` / `2` / `3` | Jump directly to that suggestion |
| `Enter` | Confirm and commit |
| `e` | Edit the selected message inline |
| `q` / `Esc` / `Ctrl+C` | Abort |

**In edit mode:**

| Key | Action |
|-----|--------|
| `Enter` | Confirm edited message |
| `Esc` | Cancel edit, return to selection |
| `←` / `→` | Move cursor |
| `Ctrl+A` / `Home` | Go to start |
| `Ctrl+E` / `End` | Go to end |
| `Backspace` | Delete character |

### Flags

```bash
ezgocommit --style gitmoji       # use gitmoji instead of conventional commits
ezgocommit --style free          # no format constraints
ezgocommit --model claude-opus-4-6  # use a different Claude model
```

## Commit styles

| Style | Example |
|-------|---------|
| `conventional` (default) | `feat(auth): add login with OAuth` |
| `gitmoji` | `✨ add login with OAuth` |
| `free` | `Add OAuth login support` |
| `custom` | Defined by you in `custom_format` |

## Documentation

- [Getting Started](docs/getting-started.md)
- [Configuration Reference](docs/configuration.md)
- [Architecture](docs/architecture.md)
- [Contributing](docs/contributing.md)

## License

MIT
