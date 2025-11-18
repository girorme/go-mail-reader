package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/BrianLeishman/go-imap"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type Credentials struct {
	Server   string
	Username string
	Password string
	Port     int
}

func getCredentials() (Credentials, error) {
	err := godotenv.Load()
	if err != nil {
		return Credentials{}, err
	}

	port, err := strconv.Atoi(os.Getenv("IMAP_PORT"))
	if err != nil {
		return Credentials{}, fmt.Errorf("[-] Invalid port: %v", err)
	}

	return Credentials{
		Server:   os.Getenv("IMAP_SERVER"),
		Port:     port,
		Username: os.Getenv("IMAP_EMAIL"),
		Password: os.Getenv("IMAP_PASSWORD"),
	}, nil
}

// ---------- Pool ----------
type IMAPPool struct {
	conns chan *imap.Dialer
	newFn func() (*imap.Dialer, error)
}

func NewIMAPPool(size int, newFn func() (*imap.Dialer, error)) (*IMAPPool, error) {
	if size <= 0 {
		return nil, fmt.Errorf("pool size must be positive, got: %d", size)
	}

	pool := &IMAPPool{
		conns: make(chan *imap.Dialer, size),
		newFn: newFn,
	}

	// Create connections concurrently for faster startup
	var wg sync.WaitGroup
	errChan := make(chan error, size)
	connChan := make(chan *imap.Dialer, size)

	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			conn, err := newFn()
			if err != nil {
				errChan <- fmt.Errorf("failed to create IMAP connection %d: %w", idx, err)
				return
			}
			if err := conn.SelectFolder("INBOX"); err != nil {
				conn.Close()
				errChan <- fmt.Errorf("failed to select folder for connection %d: %w", idx, err)
				return
			}
			connChan <- conn
		}(i)
	}

	wg.Wait()
	close(errChan)
	close(connChan)

	// Check for errors
	if len(errChan) > 0 {
		// Close any successful connections
		for conn := range connChan {
			conn.Close()
		}
		// Return first error
		return nil, <-errChan
	}

	// Add all connections to pool
	for conn := range connChan {
		pool.conns <- conn
	}

	return pool, nil
}

func (p *IMAPPool) Get() *imap.Dialer {
	return <-p.conns
}

func (p *IMAPPool) Put(conn *imap.Dialer) {
	p.conns <- conn
}

func (p *IMAPPool) Close() {
	close(p.conns)
	for c := range p.conns {
		c.Close()
	}
}

func initImap() {
	imap.Verbose = false
	imap.RetryCount = 3
}

func readMails(pool *IMAPPool, uids []int) {
	if len(uids) == 0 {
		return
	}

	// Fetch emails using one connection
	conn := pool.Get()
	emails, err := conn.GetEmails(uids...)
	pool.Put(conn)

	if err != nil {
		log.Fatalf("[-] Error getting emails: %v", err)
	}

	color.Green("[+] Processing email chunk of %d UIDs\n", len(emails))

	// Process emails concurrently with pool
	var wg sync.WaitGroup
	for _, email := range emails {
		wg.Add(1)

		go func(email *imap.Email) {
			defer wg.Done()

			color.Cyan("[+] Reading email: %s", email.Subject)

			// Get connection from pool to mark as seen
			conn := pool.Get()
			defer pool.Put(conn)

			if err := conn.MarkSeen(email.UID); err != nil {
				log.Printf("[-] Error marking email %d as seen: %v", email.UID, err)
			}
		}(email)
	}

	wg.Wait()
}

func main() {
	// Parse flags first
	chunkSize := flag.Int("chunk-size", 10, "Size of email chunks to process")
	poolSize := flag.Int("pool-size", 5, "Number of IMAP connections in the pool")
	flag.Parse()

	fmt.Println("Go mail reader")
	color.Cyan("Configuration: chunk-size=%d, pool-size=%d", *chunkSize, *poolSize)
	fmt.Println("[+] Getting envs and preparing connection")

	initImap()

	credentials, err := getCredentials()
	if err != nil {
		log.Fatalf("[-] Error getting env's: %v", err)
	}

	fmt.Printf("[+] Mail info: [%s] %s:*****\n", credentials.Server, credentials.Username)

	newConn := func() (*imap.Dialer, error) {
		im, err := imap.New(
			credentials.Username,
			credentials.Password,
			credentials.Server,
			credentials.Port,
		)
		if err != nil {
			return nil, err
		}
		return im, nil
	}

	// Create pool with configurable size
	pool, err := NewIMAPPool(*poolSize, newConn)
	if err != nil {
		log.Fatalf("[-] Error creating pool: %v", err)
	}
	defer pool.Close()

	// Use one connection to fetch unseen UIDs
	conn := pool.Get()
	fmt.Printf("[+] Getting UNSEEN UIDs\n\n")
	uids, err := conn.GetUIDs("UNSEEN")
	pool.Put(conn)

	if err != nil {
		log.Fatalf("[-] Error getting UIDs: %v", err)
	}

	if len(uids) == 0 {
		color.Yellow("[+] No unseen emails found\n")
		return
	}

	color.Green("[+] Found %d unseen emails\n", len(uids))

	// Process emails in chunks
	for uidChunk := range slices.Chunk(uids, *chunkSize) {
		readMails(pool, uidChunk)
	}

	color.Green("[+] All emails processed successfully\n")
}
