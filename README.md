# portwatch

A lightweight CLI tool that monitors open ports and alerts on unexpected changes in real time.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch && go build -o portwatch .
```

---

## Usage

Start monitoring with default settings (scans every 30 seconds):

```bash
portwatch
```

Specify a scan interval and custom port range:

```bash
portwatch --interval 10s --range 1-9999
```

Watch a specific host:

```bash
portwatch --host 192.168.1.1 --interval 5s
```

When a new port is detected or a previously open port closes, `portwatch` will print an alert to stdout:

```
[ALERT] New port detected: 8080/tcp (2024-01-15 10:32:01)
[ALERT] Port closed: 3306/tcp (2024-01-15 10:35:44)
```

### Flags

| Flag         | Default   | Description                        |
|--------------|-----------|------------------------------------|
| `--host`     | localhost | Host to monitor                    |
| `--interval` | 30s       | Interval between scans             |
| `--range`    | 1-65535   | Port range to scan                 |
| `--quiet`    | false     | Suppress output except for alerts  |

---

## License

[MIT](LICENSE)