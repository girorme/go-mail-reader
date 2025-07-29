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
		log.Fatalf("[-] Invalid port: %v", err)
	}

	return Credentials{
		Server:   os.Getenv("IMAP_SERVER"),
		Port:     port,
		Username: os.Getenv("IMAP_EMAIL"),
		Password: os.Getenv("IMAP_PASSWORD"),
	}, nil
}

func initImap() {
	imap.Verbose = false
	imap.RetryCount = 3
}

func readMails(im *imap.Dialer, uids []int) {
	emails, err := im.GetEmails(uids...)
	if err != nil {
		log.Fatalf("[-] Error getting emails: %v", err)
	}

	if len(emails) == 0 {
		fmt.Printf("[-] No emails to read on INBOX\n")
		return
	}

	var wg sync.WaitGroup

	color.Green("[+] Reading async a email chunk of %d UIDs\n", len(emails))

	for _, email := range emails {
		wg.Add(1)

		go func(email *imap.Email) {
			defer wg.Done()

			color.Cyan("[+] Reading email: %s", email.Subject)

			if err := im.MarkSeen(email.UID); err != nil {
				log.Fatalf("Error marking email as seen: %v", err)
			}
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

	im, err := imap.New(
		credentials.Username,
		credentials.Password,
		credentials.Server,
		credentials.Port,
	)
	if err != nil {
		log.Fatalf("[-] Error creating imap client: %v", err)
	}

	defer im.Close()

	fmt.Printf("[+] Selecting folder: INBOX\n")
	if err := im.SelectFolder("INBOX"); err != nil {
		log.Fatalf("[-] Error selecting folder: %v", err)
	}

	fmt.Printf("[+] Getting UNSEEN uids\n\n")
	uids, err := im.GetUIDs("UNSEEN")
	if err != nil {
		log.Fatalf("Error getting uids: %v", err)
	}

	chunkSize := flag.Int("chunk-size", 10, "Size of email chunks to process")
	flag.Parse()

	for uidChunk := range slices.Chunk(uids, *chunkSize) {
		readMails(im, uidChunk)
	}
}
