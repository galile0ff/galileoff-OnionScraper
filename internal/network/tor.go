package network

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// TorProxyAddresses kontrol edilecek portları tanımladığım kısım (Sistem Tor ve Tor Tarayıcı)
var TorProxyAddresses = []string{
	"127.0.0.1:9050", // Sistem Tor
	"127.0.0.1:9150", // Tor Browser
}

// NewTorClient trafiği Tor üzerinden yönlendiren yeni bir HTTP istemcisi oluşturur.
// Tor'un 9050 veya 9150 portunda çalışıp çalışmadığını otomatik olarak algılar.
func NewTorClient() (*http.Client, string, error) {
	var proxyAddr string
	var dialer proxy.Dialer
	var lastErr error

	// Aktif bir Tor proxy bulmaya çalış
	for _, addr := range TorProxyAddresses {

		// Bir çevirici oluştur ve bağlantıyı test et
		d, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
		if err != nil {
			lastErr = err
			continue
		}

		// Portun açık olup olmadığını kontrol eder
		conn, err := d.Dial("tcp", "check.torproject.org:80")
		if err == nil {
			conn.Close()
			proxyAddr = addr
			dialer = d
			break
		}
		lastErr = err
	}

	if proxyAddr == "" {
		return nil, "", fmt.Errorf("Tor servisi bulunamadı (9050 ve 9150 denendi). Lütfen Tor Browser'ı veya Tor servisini başlatın. Hata: %v", lastErr)
	}

	// Çeviriciyi kullanan bir transport oluşturur
	transport := &http.Transport{
		Dial: dialer.Dial,
		// Bağlantı ayarlarını optimize et
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	// Özel transport ile istemciyi döndür
	return &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}, proxyAddr, nil
}

// CheckIP o anki Tor bağlantısı üzerinden dış IP adresini sorgular
func CheckIP(client *http.Client) (string, error) {

	//Çıkış yaptığımız IP adresini kontrol etmesi için
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}
