# envlock

> Manages and validates environment variable contracts across services using a simple schema file.

---

## Installation

```bash
go install github.com/yourname/envlock@latest
```

Or build from source:

```bash
go build -o envlock .
```

---

## Usage

Define your environment contract in an `envlock.yaml` schema file:

```yaml
vars:
  DATABASE_URL:
    required: true
    description: "Primary database connection string"
  PORT:
    required: false
    default: "8080"
  API_KEY:
    required: true
    secret: true
```

Then validate your environment against the schema:

```bash
envlock validate --schema envlock.yaml
```

Check for missing or unexpected variables across services:

```bash
envlock check --schema envlock.yaml --env .env.production
```

Export a sanitized summary of your contract:

```bash
envlock export --schema envlock.yaml --format json
```

---

## Commands

| Command    | Description                                      |
|------------|--------------------------------------------------|
| `validate` | Validate current environment against schema      |
| `check`    | Compare an env file against the schema           |
| `export`   | Export schema summary in JSON or YAML format     |
| `init`     | Generate a starter `envlock.yaml` from a `.env`  |

---

## License

MIT © [yourname](https://github.com/yourname)