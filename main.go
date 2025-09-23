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
	pool := &IMAPPool{
		conns: make(chan *imap.Dialer, size),
		newFn: newFn,
	}

	for i := 0; i < size; i++ {
		conn, err := newFn()
		if err != nil {
			return nil, fmt.Errorf("failed to create imap connection: %w", err)
		}
		if err := conn.SelectFolder("INBOX"); err != nil {
			return nil, fmt.Errorf("failed to select folder: %w", err)
		}
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
	// Fetch emails using one connection (not in parallel)
	conn := pool.Get()
	emails, err := conn.GetEmails(uids...)
	pool.Put(conn)

	if err != nil {
		log.Fatalf("[-] Error getting emails: %v", err)
	}

	var wg sync.WaitGroup
	color.Green("[+] Reading async a email chunk of %d UIDs\n", len(emails))

	for _, email := range emails {
		wg.Add(1)

		go func(email *imap.Email) {
			defer wg.Done()

			color.Cyan("[+] Reading email: %s", email.Subject)

			conn := pool.Get()
			defer pool.Put(conn)

			if err := conn.MarkSeen(email.UID); err != nil {
				log.Printf("Error marking email %d as seen: %v", email.UID, err)
			}
			mu.Unlock()
		}(email)
	}

	wg.Wait()
}

func main() {
	fmt.Println("Go mail reader")
	color.Cyan("You can use -chunk-size <size> to change the size of email chunks to be read")
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

	// create pool with 5 connections
	pool, err := NewIMAPPool(5, newConn)
	if err != nil {
		log.Fatalf("[-] Error creating pool: %v", err)
	}
	defer pool.Close()

	// use one connection to fetch unseen uids
	conn := pool.Get()
	fmt.Printf("[+] Getting UNSEEN uids\n\n")
	uids, err := conn.GetUIDs("UNSEEN")
	pool.Put(conn)

	if err != nil {
		log.Fatalf("Error getting uids: %v", err)
	}

	if len(uids) == 0 {
		color.Yellow("[+] No unseen emails found\n")
		os.Exit(0)
	}

	chunkSize := flag.Int("chunk-size", 10, "Size of email chunks to process")
	flag.Parse()

	for uidChunk := range slices.Chunk(uids, *chunkSize) {
		readMails(pool, uidChunk)
	}
}
