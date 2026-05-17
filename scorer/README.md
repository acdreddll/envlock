# scorer

The `scorer` package evaluates the quality of an `envlock` schema and returns a score from 0 to 100.

## Scoring Criteria

Each environment variable can earn up to 100 weighted points:

| Criterion | Points | Condition |
|-----------|--------|-----------|
| Description | 40 | `description` field is non-empty |
| Pattern | 20 | `sensitive: true` **and** `pattern` is set |
| Default | 20 | `required: false` **and** `default` is set |
| Group | 20 | `group` field is non-empty |

The final score is the average across all variables, capped at 100.

## Usage

```go
import "github.com/user/envlock/scorer"

result := scorer.Score(mySchema)
fmt.Printf("Score: %d/100\n", result.Score)
```

## CLI

```bash
envlock score --schema envlock.yaml
envlock score --schema envlock.yaml --format json
```
