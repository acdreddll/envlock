# snapshotter

The `snapshotter` package captures point-in-time snapshots of an envlock schema and computes drift between two snapshots.

## Usage

### Take a snapshot

```go
snap := snapshotter.Take(mySchema, "v1.2.0")
```

The snapshot records the current timestamp, an optional label, and a sorted copy of all variables.

### Diff two snapshots

```go
drift := snapshotter.Diff(before, after)
for _, d := range drift {
    fmt.Printf("%s: %s\n", d.Key, d.Change)
}
```

Each `DriftEntry` contains:
- `Key` — the variable name
- `Change` — one of `added`, `removed`, `description`, `default`, `required`
- `From` / `To` — previous and new values (for field changes)

## CLI

```bash
# Capture current schema as a snapshot
envlock snapshot take --label v1 > snap_v1.json

# Compare two snapshots
envlock snapshot diff snap_v1.json snap_v2.json
envlock snapshot diff snap_v1.json snap_v2.json --format json
```

## Notes

- Snapshots do **not** capture live environment values — only schema metadata.
- Use the `pinner` package if you need to track actual runtime values.
