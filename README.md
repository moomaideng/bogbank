# Bogbank

Personal expense tracker. Architecture decisions live in [docs/arch/adr](docs/arch/adr/README.md).

## Backend

Go module `github.com/moomaideng/bogbank`. Service layout and shared libraries are documented in [services/template/README.md](services/template/README.md).

```bash
go run ./services/template serve --config services/template/config.yaml
```

`GET /livez` and `GET /readyz` listen on `:8080`.
