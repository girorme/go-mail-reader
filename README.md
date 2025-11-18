<div align="center">
  <h1>📧 go-mail-reader</h1>
  <p>A fast, efficient, and concurrent IMAP email reader built with Go</p>
  <p>
    <img src="https://img.shields.io/badge/Go-1.20%2B-blue.svg" alt="Go version">
    <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License">
  </p>
</div>

---

## 📖 Overview

`go-mail-reader` is a command-line tool designed to efficiently read and process emails from IMAP servers. It leverages Go's concurrency features with connection pooling to achieve high performance when handling large volumes of emails.

### Key Features

- ⚡ **High Performance**: Concurrent email processing with configurable connection pooling
- 🔄 **Batch Processing**: Process emails in customizable chunks for optimal memory usage
- 🔒 **Secure**: Supports TLS/SSL encrypted IMAP connections
- 🎯 **Smart Filtering**: Automatically finds and processes only unread emails
- 🛠️ **Configurable**: Flexible command-line flags for tuning performance
- 📊 **Progress Tracking**: Real-time feedback on email processing status

---

## 🎯 Use Cases

- Automated email monitoring and processing
- Bulk email reading and marking
- Email archive management
- Integration with email automation workflows
- Testing IMAP server performance

---

## 📋 Prerequisites

- **Go**: Version 1.20 or higher
- **IMAP Access**: Valid IMAP server credentials
- **TLS Support**: IMAP server with SSL/TLS support (port 993)

---

## 🛠️ Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/girorme/go-mail-reader.git
cd go-mail-reader

# Build the binary
go build -o go-mail-reader main.go

# Optional: Install globally
go install
```

### Option 2: Direct Install

```bash
go install github.com/girorme/go-mail-reader@latest
```

---

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the project root with your IMAP credentials:

```bash
# Copy the example file
cp .env.example .env

# Edit with your credentials
nano .env
```

**Required variables:**

```env
IMAP_SERVER="imap.gmail.com"     # Your IMAP server address
IMAP_PORT=993                     # IMAP port (typically 993 for SSL/TLS)
IMAP_EMAIL="your-email@gmail.com" # Your email address
IMAP_PASSWORD="your-app-password" # Your password or app-specific password
```

### Gmail Users

For Gmail, you'll need to:
1. Enable IMAP in Gmail settings
2. Use an [App Password](https://support.google.com/accounts/answer/185833) instead of your regular password
3. Set `IMAP_SERVER` to `imap.gmail.com`

---

## ⚡ Usage

### Basic Usage

```bash
# Read emails with default settings (10 emails per chunk, 5 connections)
./go-mail-reader
```

### Advanced Usage

```bash
# Process 20 emails per chunk with 10 concurrent connections
./go-mail-reader -chunk-size 20 -pool-size 10

# For smaller mailboxes (fewer connections, smaller chunks)
./go-mail-reader -chunk-size 5 -pool-size 3

# For high-volume processing (larger chunks, more connections)
./go-mail-reader -chunk-size 50 -pool-size 10
```

### Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-chunk-size` | int | 10 | Number of emails to process in each batch |
| `-pool-size` | int | 5 | Number of concurrent IMAP connections |

---

## 🚀 Performance Tips

### Optimizing for Your Use Case

1. **Large Mailboxes (1000+ emails)**
   ```bash
   ./go-mail-reader -chunk-size 50 -pool-size 10
   ```
   - Larger chunks reduce overhead
   - More connections enable parallel processing

2. **Small Mailboxes (< 100 emails)**
   ```bash
   ./go-mail-reader -chunk-size 10 -pool-size 3
   ```
   - Smaller chunks prevent over-allocation
   - Fewer connections reduce server load

3. **Slow Networks**
   ```bash
   ./go-mail-reader -chunk-size 5 -pool-size 3
   ```
   - Smaller chunks provide faster feedback
   - Fewer connections reduce timeout risks

### Connection Pool Guidelines

- **Start with defaults**: 5 connections works well for most cases
- **Monitor server limits**: Some IMAP servers limit concurrent connections
- **Balance throughput**: More connections ≠ always faster (server-dependent)
- **Typical range**: 3-10 connections for most use cases

---

## 🔍 How It Works

1. **Initialization**: Creates a pool of IMAP connections concurrently
2. **Discovery**: Fetches UIDs of all unread emails in the INBOX
3. **Batch Processing**: Divides emails into chunks for efficient processing
4. **Concurrent Execution**: Processes emails in parallel using the connection pool
5. **Mark as Read**: Automatically marks each processed email as seen

---

## 🐛 Troubleshooting

### Common Issues

**"Error getting env's"**
- Ensure `.env` file exists and contains all required variables
- Check file permissions are readable

**"Failed to create IMAP connection"**
- Verify IMAP server address and port
- Check your internet connection
- Ensure firewall allows outbound connections on port 993

**"Authentication failed"**
- Verify your email and password are correct
- For Gmail: Use an App Password, not your regular password
- Check if IMAP is enabled in your email account settings

**"Too many concurrent connections"**
- Reduce `-pool-size` value
- Some servers limit concurrent connections per account

**Slow Performance**
- Adjust `-chunk-size` and `-pool-size` based on your mailbox size
- Check network latency to IMAP server
- Verify server isn't rate-limiting your requests

---

## 📊 Example Output

```
Go mail reader
Configuration: chunk-size=10, pool-size=5
[+] Getting envs and preparing connection
[+] Mail info: [imap.gmail.com] user@gmail.com:*****
[+] Getting UNSEEN UIDs

[+] Found 25 unseen emails
[+] Processing email chunk of 10 UIDs
[+] Reading email: Important Update
[+] Reading email: Weekly Newsletter
...
[+] Processing email chunk of 10 UIDs
...
[+] Processing email chunk of 5 UIDs
...
[+] All emails processed successfully
```

---

## 🏗️ Architecture

The application uses several Go best practices:

- **Connection Pooling**: Reusable IMAP connections for efficiency
- **Concurrent Initialization**: Parallel connection setup for faster startup
- **Goroutines**: Concurrent email processing within chunks
- **Error Handling**: Comprehensive error checking and reporting
- **Clean Code**: Simple, maintainable, and well-structured

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development

```bash
# Run tests
go test -v ./...

# Format code
go fmt ./...

# Lint code
go vet ./...

# Build
go build -o go-mail-reader main.go
```

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 🙏 Acknowledgments

- Built with [go-imap](https://github.com/BrianLeishman/go-imap)
- Uses [godotenv](https://github.com/joho/godotenv) for environment management
- Colored output powered by [color](https://github.com/fatih/color)

---

## 📬 Support

If you encounter any issues or have questions:
- Open an [issue](https://github.com/girorme/go-mail-reader/issues)
- Check existing issues for solutions
- Review the troubleshooting section above

---

<div align="center">
  <p>Made with ❤️ and Go</p>
  <p>⭐ Star this repository if you find it useful!</p>
</div>

