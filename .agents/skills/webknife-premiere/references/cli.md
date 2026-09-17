# CLI

## Command naming

Commands use short, lowercase verbs: `serve`, `proxy`, `echo`, `respond`, `redirect`, `version`, `help`.

New commands should follow this pattern. Avoid compound names (`static-serve`), abbreviations (`sv`), or nouns (`server`).

## Common flags

All server commands share these flags via `addCommonServerFlags`:

| Flag | Variable | Default |
|---|---|---|
| `--listen` | `ListenAddr` | `:8080` |
| `--auth` | `Auth` | (empty) |
| `--tls-cert` | `TLSCert` | (empty) |
| `--tls-key` | `TLSKey` | (empty) |
| `--log-format` | `LogFormat` | `text` |
| `-v`, `--verbose` | `Verbose` | `false` |

Header manipulation flags are added via `addHeaderFlags`.

## Adding a new command

1. Add fields to `Config` struct
2. Create `parseNewCommand(args, stdout)` function using `flag.NewFlagSet`
3. Add `case "newcmd":` to `ParseArgs`
4. Add `case "newcmd":` to `RunWithWriter`
5. Create `buildNewHandler(cfg, bus)` function
6. Update `printUsage`

## Flag conventions

- Use `flag.NewFlagSet`, not global `flag`
- Set output to `stdout` for help text
- Use `ContinueOnError` so parse failures return errors
- Validate required flags after parsing (e.g., `--upstream` for proxy)
- Use `multiStringValue` for repeatable flags (headers)

## Error handling

- Parse errors are returned to the caller
- Invalid flag combinations are caught after parsing
- Missing required flags produce clear error messages
- The CLI prints to stderr for errors, stdout for help

## Help text

The `printUsage` function is the single source of truth for help text. Keep it updated when adding commands or flags.
