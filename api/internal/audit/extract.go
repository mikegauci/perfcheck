package audit

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

const maxBodyBytes = 2 << 20 // 2 MiB

func extractSignals(body io.Reader, base *Signals) Signals {
	sig := Signals{}
	if base != nil {
		sig = *base
	}

	z := html.NewTokenizer(io.LimitReader(body, maxBodyBytes))
	inTitle := false
	inStyle := false
	var title strings.Builder
	var style strings.Builder
	labelFors := map[string]struct{}{}
	var inputIDs []string
	var unlabeledBare int
	var imagesMissingAlt int
	var scripts, stylesheets, h1 int

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		tok := z.Token()
		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			switch tok.Data {
			case "html":
				if attr(tok, "lang") != "" {
					sig.HasLang = true
				}
			case "title":
				inTitle = true
			case "style":
				inStyle = true
			case "h1":
				h1++
			case "script":
				scripts++
				if strings.EqualFold(attr(tok, "type"), "application/ld+json") {
					sig.HasJSONLD = true
				}
			case "img":
				if !hasAttr(tok, "alt") {
					imagesMissingAlt++
				}
			case "input":
				typ := strings.ToLower(attr(tok, "type"))
				if typ == "hidden" || typ == "submit" || typ == "button" || typ == "image" {
					break
				}
				id := attr(tok, "id")
				if id == "" {
					unlabeledBare++
				} else {
					inputIDs = append(inputIDs, id)
				}
			case "label":
				if forID := attr(tok, "for"); forID != "" {
					labelFors[forID] = struct{}{}
				}
			case "meta":
				name := strings.ToLower(attr(tok, "name"))
				prop := strings.ToLower(attr(tok, "property"))
				if name == "description" {
					sig.HasMetaDescription = true
					sig.MetaDescription = attr(tok, "content")
				}
				if name == "viewport" {
					sig.HasViewport = true
				}
				if strings.HasPrefix(prop, "og:") {
					sig.HasOpenGraph = true
				}
			case "link":
				rel := strings.ToLower(attr(tok, "rel"))
				if rel == "canonical" {
					sig.HasCanonical = true
				}
				if rel == "stylesheet" {
					stylesheets++
				}
			}
		case html.EndTagToken:
			switch tok.Data {
			case "title":
				inTitle = false
			case "style":
				inStyle = false
			}
		case html.TextToken:
			if inTitle {
				title.WriteString(tok.Data)
			}
			if inStyle {
				style.WriteString(tok.Data)
			}
		}
	}

	unlabeled := unlabeledBare
	for _, id := range inputIDs {
		if _, ok := labelFors[id]; !ok {
			unlabeled++
		}
	}

	sig.Title = strings.TrimSpace(title.String())
	sig.TitleLength = len(sig.Title)
	sig.HasTitle = sig.TitleLength > 0
	sig.H1Count = h1
	sig.ImagesMissingAlt = imagesMissingAlt
	sig.InputsWithoutLabel = unlabeled
	sig.ScriptCount = scripts
	sig.StylesheetCount = stylesheets
	sig.InlineStyleBytes = style.Len()
	return sig
}

func attr(tok html.Token, key string) string {
	for _, a := range tok.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}

func hasAttr(tok html.Token, key string) bool {
	for _, a := range tok.Attr {
		if strings.EqualFold(a.Key, key) {
			return true
		}
	}
	return false
}
