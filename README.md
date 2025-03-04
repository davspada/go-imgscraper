# go-imgscraper

Downloads images from a given webpage, including dynamically loaded images and nested images.

## Requirements

Ensure you have Go installed on your system. Additionally, install the required dependencies:

```sh
go get github.com/gocolly/colly
go get github.com/chromedp/chromedp
```
Ensure you have Google Chrome installed, as ChromeDP depends on it.
## Usage

Run the program with the URL of the webpage from which you want to download images:
```sh
go run main.go <URL>
```
This will download all images from the specified URL into the images/ folder.

## Features

- Scrapes all images on the page, including dynamically loaded ones.
- Uses Colly for static images.
- Uses ChromeDP to handle JavaScript-rendered images.
- Automatically scrolls the page to load more images before scraping.
- Saves images into the images/ folder with their respective filenames.