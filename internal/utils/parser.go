package utils

import (
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks HTML içeriğinden tüm linkleri (href tag'ine bakıyor) çeker
func ExtractLinks(htmlContent string) []string {
	var links []string
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						val := strings.TrimSpace(attr.Val)
						if val != "" && !strings.HasPrefix(val, "#") && !strings.HasPrefix(val, "javascript:") {
							links = append(links, val)
						}
					}
				}
			}
		}
	}
}
