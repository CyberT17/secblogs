package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

type Blogs struct {
	Date  string
	Url   string
	Title string
	Blog  string
}

type DateBlogs struct {
	Date  string
	Blogs []Blogs
}

type HtmlValues struct {
	NumFeeds    string
	LastUpdated string
	DateBlogs   []DateBlogs
}

type RssFeed struct {
	NAME string `json:"name"`
	RSS  string `json:"rss"`
}

var tmplt *template.Template

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func main() {
	readRssFeeds()
}

func generateHtmlFile(event HtmlValues) {
	tmplt, _ = template.ParseFiles("index.gohtml")

	var f *os.File
	if _, err := os.Stat("dist/"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("dist/", os.ModeDir)
		}
	}
	f, _ = os.Create("dist/index.html")

	err := tmplt.Execute(f, event)
	check(err)

	err = f.Close()
	check(err)
}

func readRssFeeds() {
	jsonData, err := os.ReadFile("rssFeeds.json")
	check(err)
	var feeds []RssFeed

	err = json.Unmarshal(jsonData, &feeds)
	check(err)

	m := make(map[string][]Blogs)

	cutoff := time.Now().AddDate(0, 0, -7)

	feedsLength := len(feeds)
	var wg sync.WaitGroup
	wg.Add(feedsLength)

	for _, site := range feeds {
		go func(site RssFeed) {
			defer wg.Done()
			parser := gofeed.NewParser()
			parser.UserAgent = "securityblogs.xyzBot"

			feed, err := gofeed.NewParser().ParseURL(site.RSS)
			if err != nil || feed == nil {
				fmt.Print("Couldn't fetch blogs for " + site.NAME + " [")
				fmt.Print(err)
				fmt.Println("]")
				return
			}

			for _, item := range feed.Items {
				if item.PublishedParsed == nil || !item.PublishedParsed.After(cutoff) {
					continue
				}

				var dateStr = item.PublishedParsed.Format(time.DateOnly)

				singleBlog := Blogs{
					Date:  dateStr,
					Url:   item.Link,
					Title: item.Title,
					Blog:  site.NAME,
				}

				value, ok := m[dateStr]
				if !ok {
					var singleBlogList []Blogs
					m[dateStr] = append(singleBlogList, singleBlog)
				} else {
					m[dateStr] = append(value, singleBlog)
				}
			}

		}(site)
	}
	wg.Wait()
	var daten []DateBlogs
	for key, val := range m {
		singleDateBlog := DateBlogs{
			Date:  key,
			Blogs: val,
		}
		daten = append(daten, singleDateBlog)
	}
	sort.Slice(daten, func(i, j int) bool {
		return daten[i].Date > daten[j].Date
	})

	utcTime := time.Now().UTC().Format(time.RFC822)

	htmlValues := HtmlValues{
		NumFeeds:    strconv.Itoa(feedsLength),
		LastUpdated: utcTime,
		DateBlogs:   daten,
	}
	generateHtmlFile(htmlValues)
	createOpml(feeds)
}

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	Title string `xml:"title"`
}

type Body struct {
	Outline []Outline `xml:"outline"`
}

type Outline struct {
	Text string `xml:"text,attr"`
	Type string `xml:"type,attr"`
	XML  string `xml:"xmlUrl,attr"`
}

func createOpml(feeds []RssFeed) {
	var title = strconv.Itoa(len(feeds)) + " Security Blogs - securityblogs.xyz"

	opml := OPML{
		Version: "1.0",
		Head: Head{
			Title: title,
		},
		Body: Body{
			Outline: make([]Outline, len(feeds)),
		},
	}

	for i, feed := range feeds {
		opml.Body.Outline[i] = Outline{
			Text: feed.NAME,
			Type: "rss",
			XML:  feed.RSS,
		}
	}

	if _, err := os.Stat("dist/"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("dist/", os.ModeDir)
		}
	}
	opmlFile, err := os.Create("dist/securityblogs.opml")

	check(err)
	defer opmlFile.Close()

	encoder := xml.NewEncoder(opmlFile)
	encoder.Indent("", "  ")
	if err := encoder.Encode(opml); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("OPML file created successfully.")
}
