package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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

// Daktilo efekti yazma hızını buradan ayarlıyom
var TypingDelay = 2 * time.Millisecond

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
func PrintTyped(text string, delay time.Duration, color string) {
	if color != "" {
		fmt.Print(color)
	}
	for _, char := range text {
		fmt.Printf("%c", char)
		time.Sleep(delay)
	}
	if color != "" {
		fmt.Print(ColorReset)
	}
}

// ClearLine satırı temizler
func ClearLine() {
	fmt.Print("\033[2K\r")
}

// MoveCursorUp imleci yukarı taşır
func MoveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

// LiveProgress canlı ilerleme çubuğu/spinner yönetimi
type LiveProgress struct {
	stopChan chan bool
	wg       *sync.WaitGroup
	mu       sync.Mutex
	message  string
	total    int
	current  int
}

// NewLiveProgress yeni bir canlı ilerleme göstergesi oluşturur
func NewLiveProgress(message string, total int) *LiveProgress {
	return &LiveProgress{
		stopChan: make(chan bool),
		wg:       &sync.WaitGroup{},
		message:  message,
		total:    total,
		current:  0,
	}
}

// Start animasyonu başlatır
func (lp *LiveProgress) Start() {
	lp.wg.Add(1)
	go func() {
		defer lp.wg.Done()
		chars := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
		i := 0
		for {
			select {
			case <-lp.stopChan:
				lp.mu.Lock()
				ClearLine()
				lp.mu.Unlock()
				return
			default:
				lp.mu.Lock()
				percentage := 0
				if lp.total > 0 {
					percentage = (lp.current * 100) / lp.total
				}

				// Spinner + Message + Progress
				spin := string(chars[i%len(chars)])
				status := fmt.Sprintf("%s %s %s [%d/%d] %% %d", ColorCyan+spin+ColorReset, lp.message, ColorYellow, lp.current, lp.total, percentage)

				fmt.Printf("\r%s", status)
				lp.mu.Unlock()

				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

// Increment ilerlemeyi artırır
func (lp *LiveProgress) Increment() {
	lp.current++
}

// Stop animasyonu durdurur
func (lp *LiveProgress) Stop() {
	lp.stopChan <- true
	lp.wg.Wait()
	ClearLine()
}

// PrintLogWithSpinner canlı spinner çalışırken araya log girmek için kullanılır
// Spinnerı siler, logu yazar, spinnerı tekrar başlatır gibi görünür
func (lp *LiveProgress) PrintLog(logFunc func()) {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	ClearLine()
	logFunc()
}

// TypePrint formatlı metni daktilo efektiyle yazar
func TypePrint(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	PrintTyped(text, TypingDelay, "")
}

// TypePrintln satır sonu ekleyerek daktilo efekti
func TypePrintln(a ...interface{}) {
	text := fmt.Sprint(a...) + "\n"
	PrintTyped(text, TypingDelay, "")
}

// PrintBoxedTitle ASCII banner altındaki başlık kutusunu basar
func PrintBoxedTitle(title, subtitle string) {
	width := 60
	fmt.Println()

	// Çerçeve renkleri
	frameColor := ColorBlue
	titleColor := ColorBold + ColorWhite
	subColor := ColorCyan

	// Üst Kenar
	TypePrint(" %s╔%s╗%s\n", frameColor, strings.Repeat("═", width), ColorReset)

	// Başlık
	titlePad := (width - utf8.RuneCountInString(title)) / 2

	TypePrint(" %s║%s%s%s%s%s║%s\n",
		frameColor,
		strings.Repeat(" ", titlePad),
		titleColor, title,
		frameColor,
		strings.Repeat(" ", width-titlePad-utf8.RuneCountInString(title)),
		ColorReset)

	// Ara Çizgi
	TypePrint(" %s╠%s╣%s\n", frameColor, strings.Repeat("═", width), ColorReset)

	// Alt Başlık
	subPad := (width - utf8.RuneCountInString(subtitle)) / 2
	TypePrint(" %s║%s%s%s%s%s║%s\n",
		frameColor,
		strings.Repeat(" ", subPad),
		subColor, subtitle,
		frameColor,
		strings.Repeat(" ", width-subPad-utf8.RuneCountInString(subtitle)),
		ColorReset)

	// Alt Kenar
	TypePrint(" %s╚%s╝%s\n", frameColor, strings.Repeat("═", width), ColorReset)
	fmt.Println()
}

// PrintMenu seçenekleri listeler ve kullanıcı seçimini döndürür
func PrintMenu(items []string) int {
	TypePrint("\n %s%s Taranacak hedef dosyayı seçin:%s\n", ColorCyan, IconArrow, ColorReset)
	TypePrint(" %s%s%s\n", ColorBlue, strings.Repeat("┄", 40), ColorReset)

	for i, item := range items {
		TypePrint("   %s[%d]%s %s\n", ColorGreen, i+1, ColorReset, item)
	}
	TypePrint("   %s[%d]%s Manuel Dosya Yolu Gir\n", ColorYellow, len(items)+1, ColorReset)
	TypePrint("   %s[0]%s Çıkış\n", ColorRed, ColorReset)
	TypePrint(" %s%s%s\n", ColorBlue, strings.Repeat("┄", 40), ColorReset)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(" %sBirini Seç Canım >%s ", ColorCyan, ColorReset)
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

	borderColor := ColorCyan
	labelColor := ColorWhite

	// Üst Kenar
	TypePrint(" %s╔%s╗%s\n", borderColor, strings.Repeat("═", width), ColorReset)

	// Başlık
	title := "TARAMA RAPORU"
	padLeft := (width - utf8.RuneCountInString(title)) / 2
	padRight := width - padLeft - utf8.RuneCountInString(title)
	TypePrint(" %s║%s%s%s%s%s%s%s║%s\n", borderColor, ColorReset, strings.Repeat(" ", padLeft), ColorBold+ColorWhite, title, ColorReset, strings.Repeat(" ", padRight), borderColor, ColorReset)

	// Ara Çizgi
	TypePrint(" %s╠%s╣%s\n", borderColor, strings.Repeat("═", width), ColorReset)

	// Satırlar
	printReportRow("Toplam Hedef", stats.Total, ColorWhite, width, borderColor, labelColor)
	printReportRow("Başarılı", stats.Success, ColorGreen, width, borderColor, labelColor)
	printReportRow("Başarısız", stats.Failed, ColorRed, width, borderColor, labelColor)
	printReportRow("Geçen Süre", stats.Duration.Round(time.Second), ColorYellow, width, borderColor, labelColor)
	printReportRow("Klasör", stats.OutputDir, ColorPurple, width, borderColor, labelColor)

	// Alt Kenar
	TypePrint(" %s╚%s╝%s\n", borderColor, strings.Repeat("═", width), ColorReset)
	fmt.Println()
}

// AskForNewScan kullanıcıya yeni tarama yapmak isteyip istemediğini sorar
func AskForNewScan() bool {
	fmt.Println()
	TypePrint(" %s[1] Yeni Tarama Başlat%s\n", ColorGreen, ColorReset)
	TypePrint(" %s[0] Çıkış%s\n", ColorRed, ColorReset)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(" %sBirini Seç Canım >%s ", ColorCyan, ColorReset)
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

func printReportRow(label string, value interface{}, valColor string, width int, borderColor, labelColor string) {
	valStr := fmt.Sprintf("%v", value)
	valLen := utf8.RuneCountInString(valStr)

	// Padding hesabı
	// Content: " " (1) + Label (20) + " : " (3) + Value (valLen) + Padding
	// Toplam uzunluk = width olmalı
	usedWidth := 1 + 20 + 3 + valLen
	paddingLen := width - usedWidth

	if paddingLen < 0 {
		paddingLen = 0
	}

	padding := strings.Repeat(" ", paddingLen)

	// Sol boşluk, Label (sola dayalı 20), :, Değer, Sağ boşluk
	content := fmt.Sprintf(" %-20s : %s%s%s%s", label, valColor, valStr, ColorReset, padding)

	TypePrint(" %s║%s%s%s║%s\n", borderColor, labelColor, content, borderColor, ColorReset)
}

// PrintCreatedFiles oluşturulan dosyaları listeler
func PrintCreatedFiles(files []FileInfo) {
	if len(files) == 0 {
		return
	}

	TypePrint(" %s%s OLUŞTURULAN DOSYALAR:%s\n", ColorGreen, IconArrow, ColorReset)
	for _, file := range files {
		TypePrint("   %s %-30s %s(%s)%s\n", IconCheck, file.Name, ColorYellow, file.Size, ColorReset)
	}
	fmt.Println()
}

// PrintSectionHeader bölüm başlıklarını yazdırır
func PrintSectionHeader(title string) {
	fmt.Println()
	TypePrint("%s%s %s %s%s\n", ColorBlue, IconArrow, title, IconArrow, ColorReset)
	TypePrint("%s%s%s\n", ColorBlue, strings.Repeat("┄", 50), ColorReset)
}

// PrintSuccess başarılı işlem mesajını yeşil renkte yazdırır
func PrintSuccess(msg string) {
	TypePrint(" %s%s%s %s\n", ColorGreen, IconCheck, ColorReset, msg)
}

// PrintInfo bilgi mesajını mavi renkte yazdırır
func PrintInfo(msg string) {
	TypePrint(" %s%s%s %s\n", ColorBlue, IconStar, ColorReset, msg)
}

// PrintError hata mesajını kırmızı renkte yazdırır
func PrintError(msg string) {
	TypePrint(" %s%s%s %s\n", ColorRed, IconCross, ColorReset, msg)
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
	TypePrint(" %s╔%s╗%s\n", ColorRed, strings.Repeat("═", boxWidth), ColorReset)
	for _, line := range lines {
		visualLength := utf8.RuneCountInString(line)
		padding := strings.Repeat(" ", boxWidth-visualLength)
		TypePrint(" %s║ %s%s║%s\n", ColorRed, line, padding, ColorReset)
	}
	TypePrint(" %s╚%s╝%s\n", ColorRed, strings.Repeat("═", boxWidth), ColorReset)
	fmt.Println()
}

// PrintStatusLine tarama sonucunu hizalı gösterir
func PrintStatusLine(tag, item, status, detail string, isSuccess bool) {
	color := ColorRed
	icon := IconCross
	if isSuccess {
		color = ColorGreen
		icon = IconCheck
	}

	// Etiket Formatı
	tagDisplay := ""
	if tag != "" {
		tagDisplay = fmt.Sprintf("%s%-18s%s ", ColorCyan, tag, ColorReset)
	}

	TypePrint("   %s%s%-35s%s %s%s %s%s%s\n",
		tagDisplay,
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
			ClearLine()
			fmt.Printf(" %s%s%s %s... %sTAMAMLANDI%s\n", ColorGreen, IconCheck, ColorReset, action, ColorGreen, ColorReset)
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
	TypePrint(" %s%s WORKER/KÖLE (EŞZAMANLI İŞLEM) SAYISI:%s\n", ColorCyan, IconArrow, ColorReset)
	TypePrint(" %s%s%s\n", ColorBlue, strings.Repeat("┄", 40), ColorReset)
	TypePrint("   %s[1]%s 3 Worker/Köle  (Düşük - Daha yavaş, daha stabil)\n", ColorGreen, ColorReset)
	TypePrint("   %s[2]%s 5 Worker/Köle  (Orta - Dengeli)\n", ColorGreen, ColorReset)
	TypePrint("   %s[3]%s 10 Worker/Köle (Yüksek - Daha hızlı, daha fazla kaynak)\n", ColorGreen, ColorReset)
	TypePrint(" %s%s%s\n", ColorBlue, strings.Repeat("┄", 40), ColorReset)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		TypePrint(" %sBirini Seç Canım >%s ", ColorCyan, ColorReset)
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
