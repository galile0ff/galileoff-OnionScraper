package config

import (
	"bufio"
	"os"
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

// LoadTargets dosyadan URL'leri okur ve bir string dilimi döndürür
func LoadTargets(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			targets = append(targets, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}
