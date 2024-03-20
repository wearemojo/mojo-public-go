package discourseemoji

import (
	"net/url"
	"slices"
	"strings"

	"github.com/wearemojo/mojo-public-go/lib/slicefn"
	"golang.org/x/net/html"
)

func UnmungeCookedHTML(src string, baseURL *url.URL) (string, error) {
	doc, err := html.Parse(strings.NewReader(src))
	if err != nil {
		return "", err
	}

	var body *html.Node
	var fn func(*html.Node) *html.Node
	fn = func(node *html.Node) *html.Node {
		switch {
		case node.Type == html.ElementNode && node.Data == "body":
			body = node
		case node.Type == html.ElementNode && node.Data == "a" && baseURL != nil:
			if idx := slicefn.FindIndex(node.Attr, func(a html.Attribute) bool { return a.Key == "href" }); idx != -1 {
				if url, _ := url.Parse(node.Attr[idx].Val); url != nil {
					url = baseURL.ResolveReference(url)
					node.Attr[idx].Val = url.String()
				}
			}
		case node.Type == html.ElementNode && node.Data == "img":
			class, ok1 := slicefn.Find(node.Attr, func(a html.Attribute) bool { return a.Key == "class" })
			alt, ok2 := slicefn.Find(node.Attr, func(a html.Attribute) bool { return a.Key == "alt" })
			shortcode, ok3 := strings.CutPrefix(alt.Val, ":")
			shortcode, ok4 := strings.CutSuffix(shortcode, ":")
			classes := strings.Fields(class.Val)
			emoji := ShortcodeToEmoji(shortcode)

			if ok1 && ok2 && ok3 && ok4 && slices.Contains(classes, "emoji") && emoji != "" {
				return &html.Node{
					Type: html.TextNode,
					Data: emoji,
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			replacementNode := fn(child)
			if replacementNode != nil {
				node.InsertBefore(replacementNode, child)
				node.RemoveChild(child)

				child = replacementNode
			}
		}

		return nil
	}

	fn(doc)

	if body == nil {
		panic("body not found")
	}

	body.Type = html.DocumentNode

	var out strings.Builder
	if err = html.Render(&out, body); err != nil {
		return "", err
	}

	return out.String(), nil
}
