# pinner

The `pinner` package captures a **snapshot** of environment variable values and later detects **drift** — any value that has changed or disappeared since the snapshot was taken.

## Concepts

| Term | Description |
|------|-------------|
| `Snapshot` | A timestamped set of `Pin` entries created with `Capture`. |
| `Pin` | A single key/value pair recorded at a point in time. |
| `DriftReport` | The result of comparing a `Snapshot` to a live env map. |
| `DriftEntry` | A variable that has changed or gone missing. |

## Usage

```go
// Capture current values
snap := pinner.Capture(mySchema, envMap)

// Later — detect drift
report := pinner.Detect(snap, liveEnvMap)
if report.HasDrift() {
    for _, d := range report.Drifted {
        fmt.Println(d.Key, d.Pinned, "->", d.Current)
    }
}
```

## CLI

```bash
# Capture and save
envlock pin capture > snapshot.json

# Detect drift against saved snapshot
envlock pin detect snapshot.json
envlock pin detect snapshot.json --format json
```

The `detect` subcommand exits with code `1` when drift is found, making it suitable for CI gates.
