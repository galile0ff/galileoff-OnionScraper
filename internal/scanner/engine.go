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
	"galileoff-OnionScraper/internal/utils"
)

// ScanResult tarama işleminin sonucunu tutar
type ScanResult struct {
	URL        string
	StatusCode int
	Status     string
	UsedUA     string
	Error      error
	LinkCount  int
}

// StartScan bir çalışan havuzu (worker pool) ile tarama işlemini başlatır ve (başarılı, başarısız, toplam_link) sayılarını döndürür
func StartScan(targets []string, concurrency int, outputDir string) (int, int, int) {
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
		report.Log("CRITICAL", fmt.Sprintf("Tor bağlantısı başlatılamadı. Hata: %v", err))
		connectionErr = err
	} else {
		ui.PrintSuccess(fmt.Sprintf("Tor bağlantısı başarılı! Kullanılan Port: %s", proxyAddr))
		report.Log("INFO", fmt.Sprintf("Tor bağlantısı kuruldu. Port: %s", proxyAddr))
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
	totalLinks := 0

	// Sonuçları işle
	for result := range results {
		if result.Error != nil {
			failCount++
			// Hatayı log dosyasına yaz
			report.Log("FAILED", fmt.Sprintf("%s -> %v", result.URL, result.Error))

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
			totalLinks += result.LinkCount

			// Başarılı durum: HTTP Kodu ile logla
			statusText := http.StatusText(result.StatusCode)
			if statusText == "" {
				statusText = "Unknown"
			}

			// Log seviyesini belirle (200-300 SUCCESS, diğerleri WARNING)
			logLevel := "SUCCESS"
			if result.StatusCode < 200 || result.StatusCode >= 300 {
				logLevel = "WARNING"
			}

			report.Log(logLevel, fmt.Sprintf("%s -> %d %s [%s]", result.URL, result.StatusCode, statusText, result.UsedUA))

			// Başarılı mesajını göster
			ui.PrintStatusLine(result.URL, "BAŞARILI", fmt.Sprintf("(%d %s)", result.StatusCode, statusText), true)
		}
	}

	ui.PrintSectionHeader("Tarama Tamamlandı")
	return successCount, failCount, totalLinks
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

		// Başlangıç Logu
		report.Log("INFO", fmt.Sprintf("Tarama Başlatılıyor: %s (Köle Çalışmaya Başladı)", url))
		statStartTime := time.Now()

		// Request oluştur (User-Agent eklemek için)
		req, err := http.NewRequest("GET", targetURL, nil)

		// Rastgele User-Agent ve ilgili header'ları ayarla
		profile := utils.GetRandomProfile()
		if err != nil {
			// Request oluşturma hatası
			report.Log("ERROR", fmt.Sprintf("Request oluşturulamadı [%s]: %v", url, err))
			results <- ScanResult{URL: url, Status: "FAILED", UsedUA: profile.Name, Error: err}
			continue
		}

		// Header'ları ayarla
		req.Header.Set("User-Agent", profile.UserAgent)
		for k, v := range profile.Headers {
			req.Header.Set(k, v)
		}

		// İsteği gönder
		resp, err := client.Do(req)

		scanDuration := time.Since(statStartTime)

		if err != nil {
			// Bağlantı veya zaman aşımı hatası detaylı logla
			report.Log("FAILED", fmt.Sprintf("Erişim sağlanamadı [%s] (Süre: %s): %v", url, scanDuration, err))
			results <- ScanResult{URL: url, Status: "FAILED", UsedUA: profile.Name, Error: err}
			continue
		}

		statusCode := resp.StatusCode
		respSize := resp.ContentLength

		// Response Header'larını önemli olanları logla
		contentType := resp.Header.Get("Content-Type")
		server := resp.Header.Get("Server")

		report.Log("DEBUG", fmt.Sprintf("Response Alındı [%s] - Status: %d, Size: %d, Type: %s, Server: %s, Süre: %s",
			url, statusCode, respSize, contentType, server, scanDuration))

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			report.Log("ERROR", fmt.Sprintf("Response Body okunamadı [%s]: %v", url, err))
			results <- ScanResult{URL: url, Status: "FAILED", UsedUA: profile.Name, Error: err}
			continue
		}

		// HTML içeriğini kaydet
		if err := report.SaveHTML(url, string(body), outputDir); err != nil {
			report.Log("ERROR", fmt.Sprintf("%s için HTML kaydetme hatası: %v", url, err))
		} else {
			report.Log("INFO", fmt.Sprintf("HTML Kaydedildi: %s", url))
		}

		// Linkleri ayıkla ve kaydet
		links := utils.ExtractLinks(string(body))
		linkCount := len(links)

		// Her durum için links.txt dosyasına yaz
		if err := report.SaveLinks(url, links, outputDir); err != nil {
			report.Log("ERROR", fmt.Sprintf("%s için linkler kaydedilemedi: %v", url, err))
		}

		if linkCount > 0 {
			report.Log("INFO", fmt.Sprintf("%s adresinde %d adet link bulundu ve links.txt dosyasına eklendi.", url, linkCount))
			// Bulunan linkleri log dosyasına da ekle
			for _, l := range links {
				// Güvenlik: Log dosyasında da defang yapalım
				safeLink := strings.Replace(l, ".onion", "[.]onion", -1)
				report.Log("LINK", fmt.Sprintf("  -> %s", safeLink))
			}
		} else {
			report.Log("INFO", fmt.Sprintf("%s adresinde hiç link bulunamadı.", url))
		}

		// Ekran görüntüsü al (Hata olursa sadece logla, işlemi başarısız sayma)
		// Screenshot işlemi biraz zaman alacağı için köleler burada meşgul olacak
		// Ancak concurrency olduğu için diğer URL'ler işlenmeye devam ediyor
		ssStartTime := time.Now()
		if screenshotData, err := CaptureScreenshot(url, proxyAddr); err != nil {
			report.Log("FAILED", fmt.Sprintf("%s için screenshot alınamadı: %v", url, err))
		} else {
			if err := report.SaveScreenshot(url, screenshotData, outputDir); err != nil {
				report.Log("ERROR", fmt.Sprintf("%s için screenshot dosyası kaydedilemedi: %v", url, err))
			} else {
				ssDuration := time.Since(ssStartTime)
				report.Log("SUCCESS", fmt.Sprintf("%s için screenshot başarıyla kaydedildi. (Screenshot Süresi: %s)", url, ssDuration))
			}
		}

		results <- ScanResult{
			URL:        url,
			StatusCode: statusCode,
			Status:     "SUCCESS",
			UsedUA:     profile.Name,
			Error:      nil,
			LinkCount:  linkCount,
		}
	}
}
