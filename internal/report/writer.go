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

// LogHeader başlangıç başlığı yazar
func LogHeader(targetFile string, workerCount int) {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return
	}

	border := strings.Repeat("=", 60)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	header := fmt.Sprintf(`
%s
  galileoff. ONION SCRAPER - TARAMA GÜNLÜĞÜ
%s
  TARİH       : %s
  HEDEF DOSYA : %s
  ÇALIŞAN (Köle) : %d
%s
  BAŞLANGIÇ...
%s
`, border, border, timestamp, targetFile, workerCount, border, strings.Repeat("-", 60))

	logFile.WriteString(header)
}

// LogFooter kapanış özeti yazar
func LogFooter(total, success, failed int, duration time.Duration, totalSize string) {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return
	}

	border := strings.Repeat("=", 60)
	divider := strings.Repeat("-", 60)

	footer := fmt.Sprintf(`
%s
  TARAMA SONUÇ RAPORU
%s
  TOPLAM HEDEF : %d
  BAŞARILI     : %d
  BAŞARISIZ    : %d
  TOPLAM SÜRE  : %s
  VERİ BOYUTU  : %s
%s
`, divider, border, total, success, failed, duration, totalSize, border)

	logFile.WriteString(footer)
}

// Log verilen mesajı log dosyasına yazar
func Log(level, message string) {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return
	}

	// [ZAMAN] [LEVEL] Mesaj
	timestamp := time.Now().Format("15:04:05")
	entry := fmt.Sprintf("[%s] [%-7s] %s\n", timestamp, level, message)

	logFile.WriteString(entry)
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

	safeName := sanitizeFilename(url) + ".html"
	path := filepath.Join(outputDir, safeName)

	return os.WriteFile(path, []byte(content), 0644)
}

// SaveScreenshot ekran görüntüsünü belirtilen klasöre kaydeder
func SaveScreenshot(url string, data []byte, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	safeName := sanitizeFilename(url) + ".png"
	path := filepath.Join(outputDir, safeName)

	return os.WriteFile(path, data, 0644)
}

// sanitizeFilename URL'den güvenli dosya adı oluşturur
func sanitizeFilename(url string) string {
	safeName := strings.Replace(url, "http://", "", -1)
	safeName = strings.Replace(safeName, "https://", "", -1)
	safeName = strings.Replace(safeName, "/", "_", -1)
	safeName = strings.Replace(safeName, ":", "_", -1)
	return safeName
}

// Close log dosyasını kapatır
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
