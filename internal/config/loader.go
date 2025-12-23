package config

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// ListYamlFiles mevcut dizindeki taranması istenilebilecek .yaml dosyalarını listeler
func ListYamlFiles() ([]string, error) {
	var files []string
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml")) {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// LoadTargets dosyadan URL'leri okur
func LoadTargets(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)

	// Yapışık .onion adreslerini ayırmak için regex
	// .onion ibaresinden hemen sonra harf veya sayı geliyorsa araya boşluk koyar
	// Örn: site.onionjgwe... -> site.onion jgwe...
	reOnionSplit := regexp.MustCompile(`(\.onion)([a-zA-Z0-9])`)

	// .onion ibaresinden sonra / veya // varsa ve devamında yine .onionlu bir şey geliyorsa ayır
	// Örn: site.onion//baska.onion -> site.onion baska.onion
	reOnionSlashSplit := regexp.MustCompile(`(\.onion)(/{1,})([a-zA-Z0-9]+\.onion)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Güvenli hale getirilmiş (defanged) linkleri düzelt
		// [.]onion -> .onion, (.)onion -> .onion
		line = strings.ReplaceAll(line, "[.]", ".")
		line = strings.ReplaceAll(line, "(.)", ".")
		line = strings.ReplaceAll(line, "[:]", ":")
		line = strings.ReplaceAll(line, "(:)", ":")

		// Yapışık onion adreslerini ayır
		line = reOnionSplit.ReplaceAllString(line, "$1 $2")

		// Slash ile bitişik olan onion linklerini ayır
		line = reOnionSlashSplit.ReplaceAllString(line, "$1 $3")

		// Yapışık http protokollerini ayır (http://site.onionhttp://...)
		line = strings.ReplaceAll(line, "http://", " http://")
		line = strings.ReplaceAll(line, "https://", " https://")

		// Boşluklara göre parçala
		fields := strings.Fields(line)

		// Her parçayı işle ve temizle
		for _, f := range fields {
			cleaned := cleanURL(f)
			if cleaned == "" {
				continue
			}

			// Eğer .onion içeriyorsa (veya geçerli bir URL ise)
			if strings.Contains(cleaned, ".onion") || strings.HasPrefix(cleaned, "http") {
				// http:// eksikse ekle
				if !strings.HasPrefix(cleaned, "http://") && !strings.HasPrefix(cleaned, "https://") {
					cleaned = "http://" + cleaned
				}
				targets = append(targets, cleaned)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

// splitConcatenatedURLs http://a.onion/http://b.onion gibi durumları ayırır
func splitConcatenatedURLs(input string) []string {
	var result []string
	// http:// veya https:// desenlerini koruyarak stringi böl
	replaced := strings.ReplaceAll(input, "http://", "\nhttp://")
	replaced = strings.ReplaceAll(replaced, "https://", "\nhttps://")

	lines := strings.Split(replaced, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			result = append(result, cleanURL(trimmed))
		}
	}
	return result
}

// cleanURL URL sonundaki gereksiz karakterleri temizler
func cleanURL(url string) string {
	url = strings.TrimSpace(url)
	// Bazen regex yanına virgül, parantez alabilir, onları temizler
	url = strings.TrimRight(url, ",.);]'\"")
	return url
}
