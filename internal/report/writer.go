package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	logFile *os.File
	mu      sync.Mutex
)

// InitLogger log dosyasını başlat
func InitLogger(filename, outputDir string) error {
	// Log dosyasını çıktı klasörü içinde oluştur
	fullPath := filepath.Join(outputDir, filename)

	var err error
	logFile, err = os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

// Log verilen mesajı log dosyasına ve terminale yazar
func Log(message string) {
	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	if logFile != nil {
		logFile.WriteString(entry)
	}
}

// PrepareOutputDirectory belirtilen klasörü hazırlar (varsa içindekileri temizler, yoksa oluşturur)
func PrepareOutputDirectory(dirName string) error {
	// Klasör varsa içini temizle
	if _, err := os.Stat(dirName); err == nil {
		// Klasörü sil
		if err := os.RemoveAll(dirName); err != nil {
			return fmt.Errorf("klasör temizlenemedi: %v", err)
		}
	}

	// Klasörü yeniden oluştur
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("klasör oluşturulamadı: %v", err)
	}

	return nil
}

// SaveHTML kazıdığımız HTML içeriğini belirtilen klasöre kaydeder
func SaveHTML(url, content, outputDir string) error {
	// HTML dosyasını kaydetmek için klasörün varlığından emin ol
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// URL'den güvenli bir dosya adı oluştur
	// http/https protokolünü kaldır
	safeName := strings.Replace(url, "http://", "", -1)
	safeName = strings.Replace(safeName, "https://", "", -1)
	safeName = strings.Replace(safeName, "/", "_", -1) // Alt dizinleri tire yap
	safeName = strings.Replace(safeName, ":", "_", -1) // Portları tire yap

	// Uzantıyı .html yap
	if !strings.HasSuffix(safeName, ".html") {
		safeName += ".html"
	}

	path := filepath.Join(outputDir, safeName)

	return os.WriteFile(path, []byte(content), 0644)
}

// SaveScreenshot ekran görüntüsünü belirtilen klasöre kaydeder
func SaveScreenshot(url string, data []byte, outputDir string) error {
	// Klasörün varlığından emin ol
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// URL'den güvenli bir dosya adı oluştur
	safeName := strings.Replace(url, "http://", "", -1)
	safeName = strings.Replace(safeName, "https://", "", -1)
	safeName = strings.Replace(safeName, "/", "_", -1)
	safeName = strings.Replace(safeName, ":", "_", -1)

	// Uzantıyı .png yap
	if !strings.HasSuffix(safeName, ".png") {
		safeName += ".png"
	}

	path := filepath.Join(outputDir, safeName)

	return os.WriteFile(path, data, 0644)
}

// Close log dosyasını kapatır
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
