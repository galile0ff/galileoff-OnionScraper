package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"galileoff-OnionScraper/internal/config"
	"galileoff-OnionScraper/internal/network"
	"galileoff-OnionScraper/internal/report"
	"galileoff-OnionScraper/internal/scanner"
	"galileoff-OnionScraper/internal/ui"
)

func main() {
	// ASCII Banner ve başlık (Program başında bir kez)
	ui.PrintRandomBanner()
	ui.PrintBoxedTitle("galileoff. ONION SCRAPER", "Harikulade Tor Ağı Veri Kazıyıcısı")

	// IP Kontrolü
	// Kullanıcı daha menüye girmeden Tor'a bağlı mı görsün diye.
	ui.PrintInfo("IP Adresi ve Tor Bağlantısı kontrol ediliyor...")
	checkClient, _, err := network.NewTorClient()
	if err != nil {
		ui.PrintTorConnectionError(err)
	} else {
		ip, err := network.CheckIP(checkClient)
		if err != nil {
			ui.PrintError(fmt.Sprintf("IP sorgusu başarısız: %v", err))
		} else {
			ui.PrintSuccess(fmt.Sprintf("Mevcut Tor IP Adresiniz: %s", ip))
		}
	}

	for {
		// Dosya Seçimi (İnteraktif)
		targetFile := selectTargetFile()
		if targetFile == "" {
			break
		}

		// Hedefleri Yükle (Önce dosya var mı kontrol et)
		ui.PrintInfo("Hedef dosyası okunuyor: " + targetFile)
		targets, err := config.LoadTargets(targetFile)
		if err != nil {
			ui.PrintError(fmt.Sprintf("Dosya okunamadı: %v", err))
			// Hata durumunda döngü başına dön veya çıkış sor
			// Kullanıcı hatayı görüp tekrar seçim yapmak isteyebilir diye
			if !ui.AskForNewScan() {
				break
			}
			continue
		}

		// Klasör Hazırlığı
		// targets.yaml -> targets klasörü
		baseName := filepath.Base(targetFile)
		ext := filepath.Ext(baseName)
		outputDir := strings.TrimSuffix(baseName, ext)

		// Klasörü temizle/oluştur
		ui.PrintInfo(fmt.Sprintf("Çıktı klasörü hazırlanıyor: %s", outputDir))
		if err := report.PrepareOutputDirectory(outputDir); err != nil {
			ui.PrintError(fmt.Sprintf("Klasör hatası: %v", err))
			if !ui.AskForNewScan() {
				break
			}
			continue
		}

		// Loglayıcıyı Başlat
		if err := report.InitLogger("scan_report.log", outputDir); err != nil {
			ui.PrintError(fmt.Sprintf("Log dosyası oluşturulamadı: %v", err))
			if !ui.AskForNewScan() {
				break
			}
			continue
		}

		// Worker(köle) Sayısını Seç
		workerCount := ui.GetWorkerCount()

		// İstatistikleri Takip Et
		startTime := time.Now()

		// Tarayıcıyı Başlat
		successCount, failCount := scanner.StartScan(targets, workerCount, outputDir)

		report.Close() // Log dosyasını sonraki turda çakışma olmasın diye kapattık burada

		duration := time.Since(startTime)

		// Sonuç Analizi ve Raporlama
		files, totalSize := analyzeOutput(outputDir)

		stats := ui.ReportStats{
			Total:     len(targets),
			Success:   successCount,
			Failed:    failCount,
			Duration:  duration,
			DataSize:  totalSize,
			OutputDir: outputDir,
		}

		ui.PrintScanReport(stats)
		ui.PrintCreatedFiles(files)

		// Yeni tarama olacak mı?
		if !ui.AskForNewScan() {
			break
		}

		// Yeni tur başlarsa
		fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
	}
	ui.PrintTyped("İşlem Başarıyla Tamamlandı. Kendine cici bak!", 50*time.Millisecond)
}

func selectTargetFile() string {
	files, err := config.ListYamlFiles()
	if err != nil {
		ui.PrintError("Dosya listesi alınamadı.")
		return ""
	}

	choice := ui.PrintMenu(files)

	if choice == 0 {
		return ""
	}

	if choice <= len(files) {
		return files[choice-1]
	}

	// Manuel giriş
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(" %sManuel Dosya Yolu Girin >%s ", ui.ColorCyan, ui.ColorReset)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func analyzeOutput(dir string) ([]ui.FileInfo, string) {
	var files []ui.FileInfo
	var totalBytes int64

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, _ := d.Info()
			size := info.Size()
			totalBytes += size

			files = append(files, ui.FileInfo{
				Name: d.Name(),
				Size: formatSize(size),
			})
		}
		return nil
	})

	return files, formatSize(totalBytes)
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
