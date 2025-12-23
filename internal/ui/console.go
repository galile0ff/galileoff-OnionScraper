package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
)

const (
	IconCheck   = "[+]"
	IconCross   = "[-]"
	IconWarning = "[!]"
	IconArrow   = ">>"
	IconStar    = "[*]"
)

// ReportStats tarama raporu istatistiklerini tutar
type ReportStats struct {
	Total     int
	Success   int
	Failed    int
	Duration  time.Duration
	DataSize  string
	OutputDir string
}

// FileInfo oluşturulan dosya bilgilerini tutar
type FileInfo struct {
	Name string
	Size string
}

// PrintTyped metni daktilo efektiyle yazdırır
func PrintTyped(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Printf("%c", char)
		time.Sleep(delay)
	}
	fmt.Println()
}

// PrintBoxedTitle ASCII banner altındaki başlık kutusunu basar
func PrintBoxedTitle(title, subtitle string) {
	width := 60
	fmt.Println()
	// Üst çizgi: İçerik genişliği + 2 (sağ ve sol borderlar için)
	fmt.Printf(" %s%s%s\n", ColorCyan, strings.Repeat("=", width+2), ColorReset)

	// Başlığı ortala
	titlePadding := (width - utf8.RuneCountInString(title)) / 2
	fmt.Printf(" %s|%s%s%s%s|%s\n", ColorCyan, strings.Repeat(" ", titlePadding), ColorBold+title, ColorReset+ColorCyan, strings.Repeat(" ", width-titlePadding-utf8.RuneCountInString(title)), ColorReset)

	// Ara çizgi
	fmt.Printf(" %s|%s|%s\n", ColorCyan, strings.Repeat("-", width), ColorReset)

	// Alt başlığı ortala
	subPadding := (width - utf8.RuneCountInString(subtitle)) / 2
	fmt.Printf(" %s|%s%s%s|%s\n", ColorCyan, strings.Repeat(" ", subPadding), subtitle, strings.Repeat(" ", width-subPadding-utf8.RuneCountInString(subtitle)), ColorReset)

	// Alt çizgi
	fmt.Printf(" %s%s%s\n", ColorCyan, strings.Repeat("=", width+2), ColorReset)
	fmt.Println()
}

// PrintMenu seçenekleri listeler ve kullanıcı seçimini döndürür
func PrintMenu(items []string) int {
	fmt.Printf("\n %s%s Taranacak dosyayı seçin:%s\n", ColorCyan, IconArrow, ColorReset)
	fmt.Println(strings.Repeat("-", 40))

	for i, item := range items {
		fmt.Printf("   %s[%d]%s %s\n", ColorGreen, i+1, ColorReset, item)
	}
	fmt.Printf("   %s[%d]%s Manuel Dosya Yolu Gir\n", ColorGreen, len(items)+1, ColorReset)
	fmt.Printf("   %s[0]%s Çıkış\n", ColorRed, ColorReset)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(" %sNörüyoz >%s ", ColorCyan, ColorReset)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 0 && choice <= len(items)+1 {
			return choice
		}
		PrintError("Hatalı giriş! Menüdeki sayılardan birini girsene!")
	}
}

// PrintScanReport tarama sonuçlarını raporlar
func PrintScanReport(stats ReportStats) {
	fmt.Println()
	width := 50

	// Üst Çizgi
	fmt.Printf(" %s%s%s\n", ColorBlue, strings.Repeat("=", width+2), ColorReset)

	// Başlık
	title := "TARAMA RAPORU"
	padLeft := (width - utf8.RuneCountInString(title)) / 2
	padRight := width - padLeft - utf8.RuneCountInString(title)
	fmt.Printf(" %s|%s%s%s%s%s%s%s|%s\n", ColorBlue, ColorReset, strings.Repeat(" ", padLeft), ColorBold+ColorWhite, title, ColorReset, strings.Repeat(" ", padRight), ColorBlue, ColorReset)

	// Ara Çizgi
	fmt.Printf(" %s|%s|%s\n", ColorBlue, strings.Repeat("-", width), ColorReset)

	// Satırlar
	printReportRow("Toplam Hedef", stats.Total, ColorWhite, width)
	printReportRow("Başarılı", stats.Success, ColorGreen, width)
	printReportRow("Başarısız", stats.Failed, ColorRed, width)
	printReportRow("Geçen Süre", stats.Duration.Round(time.Second), ColorYellow, width)
	printReportRow("Klasör", stats.OutputDir, ColorPurple, width)

	// Alt çizgi
	fmt.Printf(" %s%s%s\n", ColorBlue, strings.Repeat("=", width+2), ColorReset)
	fmt.Println()
}

// AskForNewScan kullanıcıya yeni tarama yapmak isteyip istemediğini sorar
func AskForNewScan() bool {
	fmt.Println()
	fmt.Printf(" %s[1] Yeni Tarama Başlat%s\n", ColorGreen, ColorReset)
	fmt.Printf(" %s[0] Çıkış%s\n", ColorRed, ColorReset)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(" %sNörüyoz >%s ", ColorCyan, ColorReset)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "1" {
			return true
		}
		if input == "0" {
			return false
		}
		PrintError("Gaçaman gurtulaman, ığhhh.")
	}
}

func printReportRow(label string, value interface{}, color string, width int) {
	valStr := fmt.Sprintf("%v", value)
	valLen := utf8.RuneCountInString(valStr)

	paddingLen := width - 24 - valLen
	if paddingLen < 0 {
		paddingLen = 0
	}

	padding := strings.Repeat(" ", paddingLen)
	fmt.Printf(" %s| %s%-20s : %s%s%s%s%s|%s\n", ColorBlue, ColorReset, label, color, valStr, ColorReset, padding, ColorBlue, ColorReset)
}

// PrintCreatedFiles oluşturulan dosyaları listeler
func PrintCreatedFiles(files []FileInfo) {
	if len(files) == 0 {
		return
	}

	fmt.Printf(" %s%s OLUŞTURULAN DOSYALAR:%s\n", ColorGreen, IconArrow, ColorReset)
	for _, file := range files {
		fmt.Printf("   %s %-30s %s(%s)%s\n", IconCheck, file.Name, ColorYellow, file.Size, ColorReset)
	}
	fmt.Println()
}

// PrintSectionHeader bölüm başlıklarını yazdırır
func PrintSectionHeader(title string) {
	fmt.Println()
	fmt.Printf("%s%s %s %s%s\n", ColorBlue, IconArrow, title, IconArrow, ColorReset)
	fmt.Println(strings.Repeat("-", 50))
}

// PrintSuccess başarılı işlem mesajını yeşil renkte yazdırır
func PrintSuccess(msg string) {
	fmt.Printf(" %s%s%s %s\n", ColorGreen, IconCheck, ColorReset, msg)
}

// PrintInfo bilgi mesajını mavi renkte yazdırır
func PrintInfo(msg string) {
	fmt.Printf(" %s%s%s %s\n", ColorBlue, IconStar, ColorReset, msg)
}

// PrintError hata mesajını kırmızı renkte yazdırır
func PrintError(msg string) {
	fmt.Printf(" %s%s%s %s\n", ColorRed, IconCross, ColorReset, msg)
}

// PrintWarningBox önemli uyarıları kutu içinde gösterir
func PrintWarningBox(lines []string) {
	maxLength := 0
	for _, line := range lines {
		length := utf8.RuneCountInString(line)
		if length > maxLength {
			maxLength = length
		}
	}
	boxWidth := maxLength + 4

	fmt.Println()
	fmt.Printf("%s%s%s\n", ColorRed, strings.Repeat("=", boxWidth), ColorReset)
	for _, line := range lines {
		// Padding hesaplarken görsel uzunluğu kullan
		visualLength := utf8.RuneCountInString(line)
		padding := strings.Repeat(" ", boxWidth-4-visualLength)
		fmt.Printf("%s| %s%s |%s\n", ColorRed, line, padding, ColorReset)
	}
	fmt.Printf("%s%s%s\n", ColorRed, strings.Repeat("=", boxWidth), ColorReset)
	fmt.Println()
}

// PrintStatusLine tarama sonucunu hizalı gösterir
func PrintStatusLine(item, status, detail string, isSuccess bool) {
	color := ColorRed
	icon := IconCross
	if isSuccess {
		color = ColorGreen
		icon = IconCheck
	}

	fmt.Printf("   %s%-35s%s %s%s %s%s%s\n",
		ColorWhite, item, ColorReset,
		color, icon, status,
		ColorReset, detail)
}

// Spinner işlem sürerken yükleniyor animasyonu gösterir
func Spinner(action string, done chan bool) {
	chars := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r %s%s%s %s... %sTAMAMLANDI%s      \n", ColorGreen, IconCheck, ColorReset, action, ColorGreen, ColorReset)
			return
		default:
			fmt.Printf("\r %s%c%s %s...", ColorCyan, chars[i%len(chars)], ColorReset, action)
			time.Sleep(80 * time.Millisecond)
			i++
		}
	}
}

// PrintTorConnectionError hata penceresi gösterir
func PrintTorConnectionError(err error) {
	lines := []string{
		"TOR BAĞLANTISI BULUNAMADI",
		"Tor servisi aktif değil veya erişilemiyor.",
		"Lütfen Tor Browser'ı başlatın.",
	}

	PrintWarningBox(lines)
}

// wrapText metni belirli uzunlukta satırlara böler
func wrapText(text string, limit int) []string {
	if limit < 1 {
		return []string{text}
	}

	var lines []string
	var line strings.Builder

	words := strings.Fields(text)
	for _, word := range words {
		if line.Len()+len(word)+1 > limit {
			lines = append(lines, line.String())
			line.Reset()
		}
		if line.Len() > 0 {
			line.WriteRune(' ')
		}
		line.WriteString(word)
	}
	if line.Len() > 0 {
		lines = append(lines, line.String())
	}

	return lines
}

// GetWorkerCount kullanıcıdan worker(köle) sayısını seçmesini ister
func GetWorkerCount() int {
	fmt.Println()
	fmt.Printf(" %s%s WORKER/KÖLE (EŞZAMANLI İŞLEM) SAYISI:%s\n", ColorCyan, IconArrow, ColorReset)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("   %s[1]%s 3 Worker/Köle  (Düşük - Daha yavaş, daha stabil)\n", ColorGreen, ColorReset)
	fmt.Printf("   %s[2]%s 5 Worker/Köle  (Orta - Dengeli)\n", ColorGreen, ColorReset)
	fmt.Printf("   %s[3]%s 10 Worker/Köle (Yüksek - Daha hızlı, daha fazla kaynak)\n", ColorGreen, ColorReset)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(" %sSeçiminiz >%s ", ColorCyan, ColorReset)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "1":
			return 3
		case "2":
			return 5
		case "3":
			return 10
		default:
			PrintError("Geçersiz seçim! Lütfen 1, 2 veya 3 girin.")
		}
	}
}
