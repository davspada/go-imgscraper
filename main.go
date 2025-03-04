package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

// downloadImage downloads an image from the given URL and saves it locally
func downloadImage(imgURL, folder string) error {
	resp, err := http.Get(imgURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch image: %s", imgURL)
	}

	// Extract filename from URL
	filename := filepath.Base(strings.Split(imgURL, "?")[0])
	filePath := filepath.Join(folder, filename)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write image to file
	_, err = io.Copy(file, resp.Body)
	return err
}

func resolveURL(base, relative string) string {
	parsedBase, err := url.Parse(base)
	if err != nil {
		return relative
	}
	parsedRelative, err := url.Parse(relative)
	if err != nil {
		return relative
	}
	return parsedBase.ResolveReference(parsedRelative).String()
}

func main() {
	url := "https://www.example.com" // Change this to the target URL
	folder := "images"
	os.MkdirAll(folder, os.ModePerm)
	counter := 0

	imageURLs := make(map[string]struct{})

	// Step 1: Use Chromedp to capture dynamically loaded images with improved scrolling
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var imageNodes []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 10; i++ { // Scroll multiple times
				if err := chromedp.Run(ctx,
					chromedp.Evaluate(`window.scrollBy(0, window.innerHeight);`, nil),
					chromedp.Sleep(1*time.Second),
				); err != nil {
					return err
				}
			}
			return nil
		}),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`Array.from(document.images).map(img => img.src)`, &imageNodes),
	)
	if err != nil {
		fmt.Println("Error with ChromeDP:", err)
	}

	for _, imgURL := range imageNodes {
		if imgURL != "" {
			imageURLs[imgURL] = struct{}{}
		}
	}

	// Step 2: Use Colly to capture static images and nested ones
	c := colly.NewCollector(
		colly.AllowedDomains(strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://")),
	)

	c.OnHTML("img", func(e *colly.HTMLElement) {
		imgURL := resolveURL(url, e.Attr("src"))
		if imgURL != "" {
			imageURLs[imgURL] = struct{}{}
		}
	})

	c.Visit(url)

	// Download all images
	for imgURL := range imageURLs {
		fmt.Println("Downloading: ", imgURL)
		counter = counter + 1
		if err := downloadImage(imgURL, folder); err != nil {
			fmt.Println("Error downloading:", imgURL, err)
		}
	}

	fmt.Println("Downloaded images:", counter)
}
