package classifier

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Analyze HTML + URL analiz eder
func Analyze(htmlContent string, url string, linkCount int) Result {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return simpleAnalyze()
	}

	title := doc.Find("title").Text()
	metaDesc, _ := doc.Find("meta[name='description']").Attr("content")

	bestScore := 0
	var bestCategory *Category

	for i := range GlobalConfig.Categories {
		cat := &GlobalConfig.Categories[i]
		score := calculateScore(cat, doc, htmlContent, title, metaDesc, linkCount)

		if score > bestScore {
			bestScore = score
			bestCategory = cat
		}
	}

	// Login override: başka ciddi kategori varsa login ezilmesin
	if bestCategory != nil && bestCategory.ID == "login" {
		for i := range GlobalConfig.Categories {
			cat := &GlobalConfig.Categories[i]
			if cat.ID == "login" {
				continue
			}
			altScore := calculateScore(cat, doc, htmlContent, title, metaDesc, linkCount)
			if altScore >= bestScore-10 {
				bestScore = altScore
				bestCategory = cat
			}
		}
	}

	if bestScore < 20 || bestCategory == nil {
		return Result{
			CategoryID: "unknown",
			Tag:        "[BİLİNMEYEN]",
			Color:      "gray",
			Score:      0,
			IsUnknown:  true,
		}
	}

	return Result{
		CategoryID: bestCategory.ID,
		Tag:        bestCategory.Tag,
		Color:      bestCategory.Color,
		Score:      bestScore,
		IsUnknown:  false,
	}
}

// AnalyzeLinkContext henüz girilmemiş linkleri analiz eder
func AnalyzeLinkContext(url, anchorText string) Result {
	bestScore := 0
	var bestCategory *Category

	textLower := strings.ToLower(anchorText)
	urlLower := strings.ToLower(url)

	for i := range GlobalConfig.Categories {
		cat := &GlobalConfig.Categories[i]
		score := 0

		if strings.Contains(urlLower, cat.ID) {
			score += 5
		}

		for _, kw := range cat.Keywords.High {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score += 5
			}
		}

		for _, kw := range cat.Keywords.Medium {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score += 2
			}
		}

		for _, kw := range cat.Keywords.Exclude {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score -= 10
			}
		}

		if score > bestScore {
			bestScore = score
			bestCategory = cat
		}
	}

	if bestScore >= 5 && bestCategory != nil {
		return Result{
			CategoryID: bestCategory.ID,
			Tag:        bestCategory.Tag,
			Color:      bestCategory.Color,
			Score:      bestScore,
			IsUnknown:  false,
		}
	}

	return Result{
		CategoryID: "unknown",
		Tag:        "[?]",
		Color:      "gray",
		Score:      0,
		IsUnknown:  true,
	}
}

// calculateScore kategori skorunu hesaplar
func calculateScore(cat *Category, doc *goquery.Document, rawHTML, title, metaDesc string, linkCount int) int {
	score := 0
	lowerHTML := strings.ToLower(rawHTML)
	lowerTitle := strings.ToLower(title)
	lowerMeta := strings.ToLower(metaDesc)

	// Max link kontrolü (ceza)
	if cat.MaxLinks > 0 && linkCount > cat.MaxLinks {
		score -= 15
	}

	// Yapısal analiz
	for _, rule := range cat.StructureRules {
		if doc.Find(rule.Selector).Length() > 0 {
			score += 20
		}
	}

	// High keyword
	highHit := 0
	for _, kw := range cat.Keywords.High {
		k := strings.ToLower(kw)
		if strings.Contains(lowerHTML, k) {
			highHit++
			if highHit <= 5 {
				score += 10
			}
		}
		if strings.Contains(lowerTitle, k) || strings.Contains(lowerMeta, k) {
			score += 10
		}
	}

	// Medium keyword
	medHit := 0
	for _, kw := range cat.Keywords.Medium {
		if strings.Contains(lowerHTML, strings.ToLower(kw)) {
			medHit++
			if medHit <= 7 {
				score += 5
			}
		}
	}

	// Exclude kelimeler
	for _, kw := range cat.Keywords.Exclude {
		if strings.Contains(lowerHTML, strings.ToLower(kw)) {
			score -= 50
		}
	}

	return score
}

// fallback
func simpleAnalyze() Result {
	return Result{
		CategoryID: "unknown",
		Tag:        "[BİLİNMEYEN]",
		Color:      "gray",
		IsUnknown:  true,
	}
}
