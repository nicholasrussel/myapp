package email

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync" // Import untuk sync.WaitGroup

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var gmailService *gmail.Service

// tokenErrorResponse merepresentasikan respons error dari endpoint token
type tokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// getClient mencoba memuat token yang disimpan. Jika tidak ada,
// ia akan mendapatkan token baru menggunakan Authorization Code Flow dengan web server internal.
func getClient(config *oauth2.Config) *http.Client {

	tokFile := "/app/internal/email/token.json"

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		fmt.Println("No saved token found. Initiating Authorization Code Flow...")
		// Dapatkan token dari alur web
		tok = getTokenFromWebServer(config)
		saveToken(tokFile, tok) // Simpan token yang baru didapat
	} else {
		fmt.Println("Using saved token.")
	}
	return config.Client(context.Background(), tok)
}

// tokenFromFile membaca token dari file JSON.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// getTokenFromWebServer mendapatkan token menggunakan Authorization Code Flow
// dengan menjalankan server HTTP internal untuk menangkap pengalihan.
func getTokenFromWebServer(config *oauth2.Config) *oauth2.Token {
	// Membuat channel untuk menerima kode otorisasi
	codeChan := make(chan string)
	// Membuat WaitGroup untuk menunggu server HTTP selesai
	var wg sync.WaitGroup
	wg.Add(1)

	// Membuat server HTTP untuk menangani pengalihan OAuth
	server := &http.Server{Addr: ":8080"} // Server akan mendengarkan di port 8080

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Periksa apakah ada kode otorisasi di URL
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, "Authorization successful! You can close this tab. Token will be exchanged by the application.")
			codeChan <- code // Kirim kode otorisasi ke channel
		} else {
			fmt.Fprintf(w, "Authorization failed or code not found in URL. Please check the application logs.")
		}
		// Tutup server setelah menerima kode
		go func() {
			defer wg.Done()
			if err := server.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server shutdown error: %v", err)
			}
		}()
	})

	// Cetak URL otorisasi untuk pengguna
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("\n==================================================================================\n")
	fmt.Printf("Visit this URL on a browser (on your host machine) to authorize your application:\n")
	fmt.Printf("   %s\n", authURL)
	fmt.Printf("\nAfter authorization, your browser will redirect to http://localhost:8080.\n")
	fmt.Printf("The application will automatically capture the code.\n")
	fmt.Printf("==================================================================================\n\n")

	// Jalankan server HTTP di goroutine terpisah
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	// Tunggu kode otorisasi dari channel
	authCode := <-codeChan

	// Tunggu server HTTP benar-benar mati
	wg.Wait()

	// Tukarkan kode otorisasi dengan token
	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	fmt.Println("\nToken successfully exchanged.")
	return tok
}

// saveToken menyimpan token ke file JSON.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	// --- PERUBAHAN UTAMA DI SINI ---
	// Pastikan direktori tujuan ada sebelum menulis file
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Unable to create directory %s: %v", dir, err)
	}
	// --- AKHIR PERUBAHAN UTAMA ---

	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// InitGmailService menginisialisasi layanan Gmail API.
func InitGmailService() {
	credsPath := filepath.Join("internal", "email", "creds.json")
	b, err := os.ReadFile(credsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Gunakan ConfigFromJSON untuk membaca kredensial web application
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	gmailService, err = gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}
	fmt.Println("Gmail service initialized successfully.")
}

// SendEmail mengirim email menggunakan Gmail API.
func SendEmail(to, subject, body string) error {
	if gmailService == nil {
		return fmt.Errorf("Gmail service not initialized")
	}

	msgStr := fmt.Sprintf("From: me\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
	var message gmail.Message
	message.Raw = base64.URLEncoding.EncodeToString([]byte(msgStr))

	_, err := gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return fmt.Errorf("Unable to send email: %v", err)
	}
	return nil
}
