🇧🇷 [Português](#) | 🇺🇸 [English](#english)

---

# Primeiros Passos

## Pré-requisitos

- **Go 1.22+** — [download](https://go.dev/dl/)
- **Git** — qualquer versão recente
- **Chave de API** — escolha um dos provedores:
  - **Anthropic (Claude):** obtenha em [console.anthropic.com](https://console.anthropic.com/)
  - **Google (Gemini):** obtenha em [aistudio.google.com](https://aistudio.google.com/)

## Instalação

### A partir do código-fonte

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go build -o ezgocommit .
```

Mova o binário para algum lugar no seu `$PATH`:

```bash
sudo mv ezgocommit /usr/local/bin/
```

### Com go install

```bash
go install github.com/jeversonmisael/ez-gocommit@latest
```

### Build com tag de versão

```bash
go build -ldflags="-X github.com/jeversonmisael/ez-gocommit/cmd.Version=1.0.0" -o ezgocommit .
```

## Configurando a chave de API

A ferramenta detecta o provedor automaticamente pelo prefixo da chave: `AIzaSy*` usa Gemini, qualquer outra usa Claude.

**Com Claude (Anthropic):**

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

**Com Gemini (Google):**

```bash
export ANTHROPIC_API_KEY=AIzaSy...
```

Para torná-la permanente, adicione ao seu `~/.zshrc`, `~/.bashrc` ou equivalente.

Como alternativa, crie um arquivo de configuração em `~/.config/ezgocommit/config.toml`:

```toml
api_key = "sk-ant-..."   # ou AIzaSy... para Gemini
```

Veja [configuration.md](configuration.md) para detalhes de todas as opções.

## Primeiro uso

Faça stage de algumas mudanças e execute a ferramenta:

```bash
cd seu-projeto
git add .
ezgocommit
```

A ferramenta irá:

1. Ler seu diff staged, nome do branch e histórico de commits recentes
2. Enviar esse contexto para o Claude
3. Exibir 3 sugestões de mensagens de commit rankeadas em uma UI interativa
4. Commitar a que você escolher

## Verificando a instalação

```bash
ezgocommit version
```

## Executando sem chave de API

Se quiser testar o binário sem gastar créditos de API, você pode verificar que a ferramenta detecta a chave ausente corretamente:

```bash
ANTHROPIC_API_KEY="" ezgocommit
# Error: Anthropic API key not found.
# ...
```

## Claude Code Skill

Se você usa o [Claude Code](https://claude.ai/code), o comando `/ez-gocommit` fica disponível automaticamente após clonar o repositório. A skill está em `.claude/skills/ez-gocommit/SKILL.md`.

```text
/ez-gocommit
```

### O que a skill faz

1. Executa `git diff --cached --stat` para verificar se há arquivos staged
2. Se não houver, informa para você fazer `git add` primeiro
3. Se houver, executa `ezgocommit` — a TUI interativa abre normalmente no terminal

### Passando argumentos

Você pode passar qualquer flag do `ezgocommit` diretamente pelo slash command:

```text
/ez-gocommit --style gitmoji
/ez-gocommit --style free
/ez-gocommit --model claude-opus-4-6
```

### Modelo de permissões

A skill declara `allowed-tools: Bash(git *), Bash(ezgocommit *)`, o que limita os comandos que o Claude Code pode executar apenas a operações Git e ao próprio binário. Nenhum outro comando do sistema é permitido.

Além disso, `disable-model-invocation: true` garante que o Claude não tenta gerar a mensagem de commit por conta própria — a análise fica inteiramente com o `ezgocommit`.

> Você sempre invoca a skill explicitamente com `/ez-gocommit`; ela nunca roda de forma automática.

## Próximos passos

- [Referência de Configuração](configuration.md) — personalizar modelo, estilo, limites de diff
- [Arquitetura](architecture.md) — entender como o código está estruturado
- [Contribuindo](contributing.md) — como adicionar funcionalidades ou corrigir bugs

---

<a id="english"></a>

🇧🇷 [Português](#) | 🇺🇸 [English](#english)

# Getting Started

## Prerequisites

- **Go 1.22+** — [download](https://go.dev/dl/)
- **Git** — any recent version
- **API key** — choose a provider:
  - **Anthropic (Claude):** get one at [console.anthropic.com](https://console.anthropic.com/)
  - **Google (Gemini):** get one at [aistudio.google.com](https://aistudio.google.com/)

## Installation

### From source

```bash
git clone https://github.com/jeversonmisael/ez-gocommit
cd ez-gocommit
go build -o ezgocommit .
```

Move the binary to somewhere in your `$PATH`:

```bash
sudo mv ezgocommit /usr/local/bin/
```

### With go install

```bash
go install github.com/jeversonmisael/ez-gocommit@latest
```

### Build with version tag

```bash
go build -ldflags="-X github.com/jeversonmisael/ez-gocommit/cmd.Version=1.0.0" -o ezgocommit .
```

## Setting up the API key

The tool detects the provider automatically from the key prefix: `AIzaSy*` uses Gemini, anything else uses Claude.

**With Claude (Anthropic):**

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

**With Gemini (Google):**

```bash
export ANTHROPIC_API_KEY=AIzaSy...
```

To make it permanent, add that line to your `~/.zshrc`, `~/.bashrc`, or equivalent.

Alternatively, create a config file at `~/.config/ezgocommit/config.toml`:

```toml
api_key = "sk-ant-..."   # or AIzaSy... for Gemini
```

See [configuration.md](configuration.md) for details on all options.

## First use

Stage some changes and run the tool:

```bash
cd your-project
git add .
ezgocommit
```

The tool will:

1. Read your staged diff, branch name, and recent commit history
2. Send that context to Claude
3. Display 3 ranked commit message suggestions in an interactive UI
4. Commit the one you choose

## Verifying the install

```bash
ezgocommit version
```

## Running without an API key

If you want to test the binary without spending API credits, you can check that the tool detects the missing key correctly:

```bash
ANTHROPIC_API_KEY="" ezgocommit
# Error: Anthropic API key not found.
# ...
```

## Claude Code Skill

If you use [Claude Code](https://claude.ai/code), the `/ez-gocommit` command is available automatically after cloning the repo. The skill is at `.claude/skills/ez-gocommit/SKILL.md`.

```text
/ez-gocommit
```

### What the skill does

1. Runs `git diff --cached --stat` to check for staged files
2. If none are found, it tells you to `git add` first
3. If there are staged changes, it runs `ezgocommit` — the interactive TUI opens normally in the terminal

### Passing arguments

You can pass any `ezgocommit` flag directly through the slash command:

```text
/ez-gocommit --style gitmoji
/ez-gocommit --style free
/ez-gocommit --model claude-opus-4-6
```

### Permission model

The skill declares `allowed-tools: Bash(git *), Bash(ezgocommit *)`, which limits what Claude Code can execute to Git operations and the binary itself. No other system commands are permitted.

Additionally, `disable-model-invocation: true` ensures Claude does not try to generate the commit message on its own — all analysis stays within `ezgocommit`.

> You always invoke the skill explicitly with `/ez-gocommit`; it never runs automatically.

## Next steps

- [Configuration Reference](configuration.md) — customize model, style, diff limits
- [Architecture](architecture.md) — understand how the codebase is structured
- [Contributing](contributing.md) — how to add features or fix bugs
