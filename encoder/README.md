# encoder

The `encoder` package encodes environment variable values into a chosen format (base64, hex, or none) based on the envlock schema.

## Usage

```go
import "github.com/envlock/encoder"

results, err := encoder.Encode(s, env, encoder.Options{
    Format:        encoder.FormatBase64,
    SensitiveOnly: true,
})
```

## Options

| Field           | Type     | Description                                   |
|-----------------|----------|-----------------------------------------------|
| `Format`        | `Format` | `base64`, `hex`, or `none` (passthrough)       |
| `SensitiveOnly` | `bool`   | When true, only encodes vars marked sensitive  |

## Result Fields

| Field      | Description                                      |
|------------|--------------------------------------------------|
| `Key`      | Variable key name                                |
| `Original` | Original plaintext value                         |
| `Encoded`  | Encoded value (or original if skipped/none)      |
| `Format`   | Format used                                      |
| `Skipped`  | True if var was absent from env or not targeted  |

## CLI

```
envlock encode --schema envlock.yaml --format base64 --sensitive-only
envlock encode --format hex --output json
```
