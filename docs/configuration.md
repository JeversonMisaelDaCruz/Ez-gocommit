🇧🇷 [Português](#) | 🇺🇸 [English](#english)

---

# Referência de Configuração

O Ez-gocommit é configurado através de variáveis de ambiente e/ou arquivos TOML. Variáveis de ambiente sempre têm prioridade.

## Chave de API

A chave de API é a única configuração obrigatória.

| Método | Valor |
|--------|-------|
| Variável de ambiente | `ANTHROPIC_API_KEY=sk-ant-...` |
| Campo no arquivo de config | `api_key = "sk-ant-..."` |

A variável de ambiente tem precedência sobre o arquivo de configuração.

## Locais do arquivo de configuração

A ferramenta lê arquivos de configuração nesta ordem (o último vence para cada chave):

1. `~/.config/ezgocommit/config.toml` — configuração global do usuário
2. `.ezgocommit.toml` no diretório atual — configuração no nível do projeto

A configuração do projeto substitui a configuração global. Isso permite definir um estilo de commit diferente por repositório.

## Todas as opções

| Campo | Tipo | Padrão | Descrição |
|-------|------|--------|-----------|
| `api_key` | string | — | Chave de API Anthropic (preferir variável de ambiente) |
| `model` | string | `claude-sonnet-4-6` | Modelo Claude a usar |
| `commit_style` | string | `conventional` | Formato da mensagem: `conventional`, `gitmoji`, `free`, `custom` |
| `custom_format` | string | — | Descreva seu formato quando `commit_style = "custom"` |
| `language` | string | `en` | Idioma das mensagens geradas |
| `max_diff_lines` | int | `500` | Máximo de linhas de diff enviadas para a IA (evita prompts enormes) |

## Exemplo de arquivo de configuração

```toml
# ~/.config/ezgocommit/config.toml

api_key        = "sk-ant-..."
model          = "claude-sonnet-4-6"
commit_style   = "conventional"
language       = "pt"
max_diff_lines = 500
```

## Sobrescrita por projeto

Coloque `.ezgocommit.toml` na raiz de um repositório:

```toml
# .ezgocommit.toml — este repositório usa gitmoji
commit_style = "gitmoji"
```

## Flags de linha de comando

Flags substituem tanto os arquivos de configuração quanto as variáveis de ambiente para aquela execução:

```bash
ezgocommit --style gitmoji
ezgocommit --model claude-opus-4-6
```

| Flag | Substitui |
|------|-----------|
| `--style` | `commit_style` |
| `--model` | `model` |
| `--config` | caminho do arquivo de config (reservado, ainda não implementado) |

## Estilos de commit

### `conventional` (padrão)

Segue a especificação [Conventional Commits](https://www.conventionalcommits.org/).

```
feat(scope): short description
fix: correct null pointer in auth handler
chore(deps): update go modules
```

Tipos comuns: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `style`, `perf`, `ci`, `build`

### `gitmoji`

Usa prefixos de emoji do [gitmoji](https://gitmoji.dev/).

```
✨ add OAuth login support
🐛 fix race condition in token refresh
♻️ refactor user repository layer
```

### `free`

Sem restrições de formato. A IA escreve o que achar que melhor descreve a mudança.

```
Add OAuth login support
Fix race condition when refreshing tokens
Clean up user repository layer
```

### `custom`

Defina `commit_style = "custom"` e descreva seu formato em `custom_format`. A IA seguirá exatamente.

```toml
commit_style   = "custom"
custom_format  = "JIRA-XXXX | tipo: descrição curta"
```

## Modelos Claude disponíveis

| Modelo | ID | Notas |
|--------|----|-------|
| Sonnet 4.6 (padrão) | `claude-sonnet-4-6` | Melhor equilíbrio entre qualidade e velocidade |
| Opus 4.6 | `claude-opus-4-6` | Maior qualidade, mais lento |
| Haiku 4.5 | `claude-haiku-4-5-20251001` | Mais rápido, mais econômico |

---

<a id="english"></a>

🇧🇷 [Português](#) | 🇺🇸 [English](#english)

# Configuration Reference

Ez-gocommit is configured through environment variables and/or TOML files. Environment variables always take priority.

## API key

The API key is the only required setting.

| Method | Value |
|--------|-------|
| Environment variable | `ANTHROPIC_API_KEY=sk-ant-...` |
| Config file field | `api_key = "sk-ant-..."` |

The environment variable takes precedence over the config file.

## Config file locations

The tool reads config files in this order (last one wins for each key):

1. `~/.config/ezgocommit/config.toml` — global user config
2. `.ezgocommit.toml` in the current directory — project-level config

Project-level config overrides global config. This lets you set a different commit style per repository.

## All options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `api_key` | string | — | Anthropic API key (prefer env var) |
| `model` | string | `claude-sonnet-4-6` | Claude model to use |
| `commit_style` | string | `conventional` | Message format: `conventional`, `gitmoji`, `free`, `custom` |
| `custom_format` | string | — | Describe your format when `commit_style = "custom"` |
| `language` | string | `en` | Language for generated messages |
| `max_diff_lines` | int | `500` | Max diff lines sent to the AI (prevents huge prompts) |

## Example config file

```toml
# ~/.config/ezgocommit/config.toml

api_key        = "sk-ant-..."
model          = "claude-sonnet-4-6"
commit_style   = "conventional"
language       = "en"
max_diff_lines = 500
```

## Per-project override

Place `.ezgocommit.toml` in the root of a repository:

```toml
# .ezgocommit.toml — this repo uses gitmoji
commit_style = "gitmoji"
```

## CLI flags

Flags override both config files and environment variables for that single run:

```bash
ezgocommit --style gitmoji
ezgocommit --model claude-opus-4-6
```

| Flag | Overrides |
|------|-----------|
| `--style` | `commit_style` |
| `--model` | `model` |
| `--config` | config file path (reserved, not yet implemented) |

## Commit styles

### `conventional` (default)

Follows the [Conventional Commits](https://www.conventionalcommits.org/) specification.

```
feat(scope): short description
fix: correct null pointer in auth handler
chore(deps): update go modules
```

Common types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `style`, `perf`, `ci`, `build`

### `gitmoji`

Uses [gitmoji](https://gitmoji.dev/) emoji prefixes.

```
✨ add OAuth login support
🐛 fix race condition in token refresh
♻️ refactor user repository layer
```

### `free`

No format constraints. The AI writes what it thinks best describes the change.

```
Add OAuth login support
Fix race condition when refreshing tokens
Clean up user repository layer
```

### `custom`

Set `commit_style = "custom"` and describe your format in `custom_format`. The AI will follow it exactly.

```toml
commit_style   = "custom"
custom_format  = "JIRA-XXXX | type: short description"
```

## Available Claude models

| Model | ID | Notes |
|-------|----|-------|
| Sonnet 4.6 (default) | `claude-sonnet-4-6` | Best balance of quality and speed |
| Opus 4.6 | `claude-opus-4-6` | Highest quality, slower |
| Haiku 4.5 | `claude-haiku-4-5-20251001` | Fastest, most economical |
