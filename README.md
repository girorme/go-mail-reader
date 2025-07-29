<div align="center">
  <h1>📧 go-mail-reader</h1>
  <p>A simple and efficient command-line tool to read emails using Go.</p>
  <img src="https://img.shields.io/badge/Go-1.20%2B-blue.svg" alt="Go version">
</div>

---

## 🚀 Features

- Read emails from IMAP servers
- Easy to configure and extend

---

## 🛠️ Installation

```bash
git clone https://github.com/girorme/go-mail-reader.git
cd go-mail-reader
go build -o go-mail-reader main.go
```

---

## ⚡ Usage

```bash
./go-mail-reader -chunk-size <size>
```

#### Common Flags

- `--server`   : IMAP server address
- `--user`     : Email address
- `--password` : Email password or app password
- `--folder`   : Mailbox folder (default: INBOX)
- `--since`    : Fetch emails since date (YYYY-MM-DD)
- `--format`   : Output format (`text` or `json`)

---

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---
