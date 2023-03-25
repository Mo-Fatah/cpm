package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	Url        string
	htmlBody   io.Reader
	hashedJobs string
	JobLinks   []string
}

func NewCrawler(url string) Crawler {
	return Crawler{
		Url:      url,
		htmlBody: nil,
		JobLinks: nil,
	}
}

func (c *Crawler) Crawl() (io.Reader, error) {
	response, err := http.Get(c.Url)
	if err != nil {
		return nil, err
	}
	c.htmlBody = response.Body
	return response.Body, nil
}

// Parse the given html page and extract the relevant jobs links.
// It calls Crawl if it hasn't been called yet.
func (c *Crawler) Parse() error {
	if c.htmlBody == nil {
		if _, err := c.Crawl(); err != nil {
			return err
		}
	}

	doc, err := goquery.NewDocumentFromReader(c.htmlBody)
	if err != nil {
		return err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
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
	if len(elements) < 2 || !isVariantOfJob(elements[0]) {
		return false
	}
	return true
}

// some websites call the resource job(s) and others career(s)
func isVariantOfJob(str string) bool {
	if str == "job" || str == "jobs" || str == "career" || str == "careers" {
		return true
	}
	return false
}
