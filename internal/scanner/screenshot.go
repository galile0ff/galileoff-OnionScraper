package scanner

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// CaptureScreenshot belirtilen URL'in ekran görüntüsünü alır
func CaptureScreenshot(url, proxyAddr string) ([]byte, error) {
	// Edge tarayıcısının yolunu bulmaya çalış
	edgePaths := []string{
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
		`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
	}

	var execPath string
	for _, path := range edgePaths {
		if _, err := os.Stat(path); err == nil {
			execPath = path
			break
		}
	}

	// Proxy ayarlarını yap
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("socks5://"+proxyAddr),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		// WebRTC ile IP sızıntısı önleme
		chromedp.Flag("disable-webrtc", true),
		chromedp.Flag("force-webrtc-ip-handling-policy", "disable_non_proxied_udp"),
		// DNS sızıntı önleme
		chromedp.Flag("host-resolver-rules", "MAP * ~NOTFOUND , EXCLUDE 127.0.0.1"),
		// Gereksiz servisleri kapatma
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.WindowSize(1280, 1024),
	)

	// Eğer Edge bulunduysa onu kullan
	if execPath != "" {
		opts = append(opts, chromedp.ExecPath(execPath))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Zaman aşımı bağlamı oluştur (30 saniye iyi gibi)
	ctx, cancel := context.WithTimeout(allocCtx, 30*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var buf []byte

	// Eğer eksikse http:// önekini ekle
	targetURL := url
	if len(url) > 0 && url[:4] != "http" {
		targetURL = "http://" + url
	}

	// Görevleri çalıştır
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(2*time.Second), // Sayfanın tam yüklenmesi için sabır
		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		return nil, fmt.Errorf("ekran görüntüsü alınamadı: %v", err)
	}

	return buf, nil
}
