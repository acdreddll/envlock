# classifier

The `classifier` package automatically assigns a **category** to each environment variable in a schema based on its metadata, naming conventions, and tags.

## Categories

| Category       | Description                                              |
|----------------|----------------------------------------------------------|
| `secret`       | Sensitive credentials, tokens, keys, passwords           |
| `feature_flag` | Feature toggles and boolean enable/disable switches      |
| `infra`        | Infrastructure concerns: hosts, ports, URLs, regions     |
| `config`       | General application configuration (default)              |
| `unknown`      | Could not be classified                                  |

## Classification Rules (in priority order)

1. **Sensitive flag** — if `sensitive: true` → `secret`
2. **Tags** — if any tag matches `secret`, `credential`, `auth`, `feature`, `flag`, `toggle`, `infra`, etc.
3. **Key prefix** — e.g. `TOKEN_`, `ENABLE_`, `HOST_`
4. **Key suffix** — e.g. `_KEY`, `_TOKEN`, `_URL`, `_PORT`
5. **Default** → `config`

## Usage

```go
import "github.com/envlock/classifier"

results := classifier.Classify(mySchema)
for _, r := range results {
    fmt.Printf("%s → %s (%s)\n", r.Key, r.Category, r.Reason)
}
```

## CLI

```bash
envlock classify --format text
envlock classify --format json
```
