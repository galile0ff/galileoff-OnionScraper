package scanner

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"galileoff-OnionScraper/internal/network"
	"galileoff-OnionScraper/internal/report"
	"galileoff-OnionScraper/internal/ui"
)

// ScanResult tarama işleminin sonucunu tutar
type ScanResult struct {
	URL    string
	Status string
	Error  error
}

// StartScan bir çalışan havuzu (worker pool) ile tarama işlemini başlatır ve (başarılı, başarısız) sayılarını döndürür
func StartScan(targets []string, concurrency int, outputDir string) (int, int) {
	client, proxyAddr, err := network.NewTorClient()

	// Tor bağlantı durumu kontrolü
	var connectionErr error
	if err != nil {
		ui.PrintWarningBox([]string{
			"TOR BAĞLANTISI BULUNAMADI",
			"Tor servisi aktif değil.",
			".onion sitelerine erişim sağlanamayacak.",
			"Lütfen Tor Browser'ı başlatıp tekrar deneyin.",
		})

		// Hatayı olduysa terminale ve log dosyasına logla
		report.Log(fmt.Sprintf("[CRITICAL_ERROR] Tor bağlantısı başlatılamadı. Hata: %v", err))
		connectionErr = err
	} else {
		ui.PrintSuccess(fmt.Sprintf("Tor bağlantısı başarılı! Kullanılan Port: %s", proxyAddr))
		report.Log(fmt.Sprintf("[INFO] Tor bağlantısı kuruldu. Port: %s", proxyAddr))
		ui.PrintInfo("Gizlilik Modu: Tor Browser İmzası (User-Agent) Aktif")
	}

	tasks := make(chan string, len(targets))
	results := make(chan ScanResult, len(targets))
	var wg sync.WaitGroup

	ui.PrintSectionHeader(fmt.Sprintf("Tarama Başlatılıyor (%d Hedef)", len(targets)))

	// İşçileri (workers/köle) başlat
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(client, proxyAddr, tasks, results, &wg, connectionErr, outputDir)
	}

	// Görevleri gönder
	for _, target := range targets {
		tasks <- target
	}
	close(tasks)

	// Tamamlanmayı ayrı bir goroutine'de bekle
	go func() {
		wg.Wait()
		close(results)
	}()

	successCount := 0
	failCount := 0

	// Sonuçları işle
	for result := range results {
		if result.Error != nil {
			failCount++
			// Hatayı log dosyasına yaz
			report.Log(fmt.Sprintf("[FAILED] %s -> %v", result.URL, result.Error))

			// Kullanıcı için basit bir mesaj
			var detailMsg string
			if strings.Contains(result.Error.Error(), "TOR_CONNECTION_MISSING") {
				detailMsg = "(TOR Servisi Yok)"
			} else {
				detailMsg = "(Erişim Hatası)"
			}
			ui.PrintStatusLine(result.URL, "BAŞARISIZ", detailMsg, false)
		} else {
			successCount++
			// Başarılı durum
			report.Log(fmt.Sprintf("[SUCCESS] %s -> OK", result.URL))

			// Başarılı mesajını göster
			ui.PrintStatusLine(result.URL, "BAŞARILI", "(Veri Çekildi)", true)
		}
	}

	ui.PrintSectionHeader("Tarama Tamamlandı")
	return successCount, failCount
}

func worker(client *http.Client, proxyAddr string, tasks <-chan string, results chan<- ScanResult, wg *sync.WaitGroup, connectionErr error, outputDir string) {
	defer wg.Done()
	for url := range tasks {
		// Eğer Tor bağlantısı baştan yoksa direkt hata dön
		if connectionErr != nil {
			results <- ScanResult{
				URL:    url,
				Status: "FAILED",
				Error:  fmt.Errorf("TOR_CONNECTION_MISSING: Tor bağlantısı olmadığı için erişilemedi. (%v)", connectionErr),
			}
			// UI'ın tıkanmaması için (logların anında basılmasını önlemek için) delay
			time.Sleep(50 * time.Millisecond)
			continue
		}

		// Eğer eksikse http:// önekini ekle
		targetURL := url
		if len(url) > 0 && url[:4] != "http" {
			targetURL = "http://" + url
		}

		// Request oluştur (User-Agent eklemek için)
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			results <- ScanResult{URL: url, Status: "FAILED", Error: err}
			continue
		}

		// Tor Browser User-Agent'ı ekle (Gizlilik için standart Tor UA kullandım)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0")

		resp, err := client.Do(req)
		if err != nil {
			results <- ScanResult{URL: url, Status: "FAILED", Error: err}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			results <- ScanResult{URL: url, Status: "FAILED", Error: err}
			continue
		}

		// HTML içeriğini kaydet
		if err := report.SaveHTML(url, string(body), outputDir); err != nil {
			report.Log(fmt.Sprintf("%s için HTML kaydetme hatası: %v", url, err))
		}

		// Ekran görüntüsü al (Hata olursa sadece logla, işlemi başarısız sayma)
		// Screenshot işlemi biraz zaman alacağı için köleler burada meşgul olacak
		// Ancak concurrency olduğu için diğer URL'ler işlenmeye devam ediyor
		if screenshotData, err := CaptureScreenshot(url, proxyAddr); err != nil {
			report.Log(fmt.Sprintf("[FAILED] %s için screenshot alınamadı: %v", url, err))
		} else {
			if err := report.SaveScreenshot(url, screenshotData, outputDir); err != nil {
				report.Log(fmt.Sprintf("[FAILED] %s için screenshot dosyası kaydedilemedi: %v", url, err))
			} else {
				report.Log(fmt.Sprintf("[SUCCESS] %s için screenshot kaydedildi.", url))
			}
		}

		results <- ScanResult{URL: url, Status: "SUCCESS", Error: nil}
	}
}
