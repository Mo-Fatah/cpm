package internal

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Crawler struct {
	Url        string
	queryDoc   *goquery.Document
	hashedJobs string
	JobLinks   []string
}

func NewCrawler(url string) Crawler {
	return Crawler{
		Url:      url,
		queryDoc: nil,
		JobLinks: nil,
	}
}

func (c *Crawler) Crawl() error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(c.Url),
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)

	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return err
	}
	c.queryDoc = doc
	return nil
}

// Parse the given html page and extract the relevant jobs links.
// It calls Crawl if it hasn't been called yet.
func (c *Crawler) Parse() error {
	if c.queryDoc == nil {
		if err := c.Crawl(); err != nil {
			return err
		}
	}

	c.queryDoc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			// Ignore mailto: and tel: links
			if strings.HasPrefix(link, "mailto:") && strings.HasPrefix(link, "tel:") {
				return
			}
			if isJobLink(link) {
				c.JobLinks = append(c.JobLinks, link)
			}
		}
	})
	return nil
}

func (c *Crawler) GetJobLinks() ([]string, error) {
	if c.JobLinks == nil {
		if err := c.Parse(); err != nil {
			return nil, err
		}
	}
	return c.JobLinks, nil
}

// should be called after Crawl() and Parse()
func (c *Crawler) GetHash() string {
	if len(c.hashedJobs) == 0 {
		h := sha1.New()
		for _, job := range c.JobLinks {
			h.Write([]byte(job))
		}
		c.hashedJobs = hex.EncodeToString((h.Sum(nil)))
	}
	return c.hashedJobs
}

// extract the main name of the board from the url
func (c *Crawler) GetBoardName() string {
	elements := strings.Split(c.Url, "/")
	return elements[2]
}

// check that link at the form <resource>/<id>
func isJobLink(link string) bool {
	link = strings.Trim(link, "/")
	elements := strings.Split(link, "/")
	if len(elements) < 2 || !isVariantOfJob(link) {
		return false
	}
	return true
}

// some websites call the resource job(s) and others career(s)
func isVariantOfJob(link string) bool {

	if strings.Contains(link, "job") || strings.Contains(link, "career") || strings.Contains(link, "jobs") || strings.Contains(link, "careers") {
		return true
	}
	return false
}
