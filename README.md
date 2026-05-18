# logslice

Stream and filter structured JSON logs from multiple sources with a query DSL and human-readable output.

---

## Installation

```bash
go install github.com/yourname/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/logslice.git && cd logslice && go build -o logslice .
```

---

## Usage

```bash
# Filter logs from a file
logslice --file app.log --query 'level == "error"'

# Stream logs from multiple sources
logslice --file app.log --file worker.log --follow

# Filter by field value and format output
logslice --file app.log --query 'status >= 500 AND service == "api"' --format pretty

# Pipe from stdin
kubectl logs -f my-pod | logslice --query 'level == "warn" OR level == "error"'
```

### Query DSL

| Operator | Example |
|----------|---------|
| `==`     | `level == "error"` |
| `!=`     | `service != "worker"` |
| `>=`     | `latency_ms >= 200` |
| `AND / OR` | `level == "error" AND env == "prod"` |

### Output Formats

- `pretty` — Human-readable, colorized output (default)
- `json` — Raw JSON passthrough
- `compact` — Single-line key=value pairs

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

MIT © yourname