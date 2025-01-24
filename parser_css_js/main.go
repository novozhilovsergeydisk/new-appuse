package main
// go get golang.org/x/net/html
import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

// downloadFile downloads the file from the given URL and saves it to the specified filename.
func downloadFile(fileURL, filename string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", fileURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save content to file %s: %w", filename, err)
	}

	return nil
}

// extractLinks parses the HTML and extracts all CSS and JS links.
func extractLinks(htmlFile string) ([]string, error) {
	file, err := os.Open(htmlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open HTML file: %w", err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var links []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "link" {
				for _, attr := range n.Attr {
					if attr.Key == "href" && (strings.HasSuffix(attr.Val, ".css") || strings.Contains(attr.Val, ".css?")) {
						links = append(links, attr.Val)
					}
				}
			} else if n.Data == "script" {
				for _, attr := range n.Attr {
					if attr.Key == "src" && (strings.HasSuffix(attr.Val, ".js") || strings.Contains(attr.Val, ".js?")) {
						links = append(links, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return links, nil
}

// saveAssets downloads and saves each asset to a file.
func saveAssets(links []string) error {
	for _, link := range links {
		parsedURL, err := url.Parse(link)
		if err != nil {
			return fmt.Errorf("invalid URL %s: %w", link, err)
		}

		filename := path.Base(parsedURL.Path)
		fmt.Printf("Downloading %s to %s...\n", link, filename)
		if err := downloadFile(link, filename); err != nil {
			return fmt.Errorf("failed to download %s: %w", link, err)
		}
	}
	return nil
}

func main() {
	// Path to the HTML file to parse
	const htmlFile = "index.html"

	// Extract CSS and JS links
	links, err := extractLinks(htmlFile)
	if err != nil {
		fmt.Printf("Error extracting links: %v\n", err)
		os.Exit(1)
	}

	// Save assets to the current directory
	if err := saveAssets(links); err != nil {
		fmt.Printf("Error saving assets: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Assets downloaded successfully.")
}
