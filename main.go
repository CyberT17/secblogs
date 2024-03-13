package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
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
	Date string
	Blogs []Blogs
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

func generateHtmlFile(event []DateBlogs) {
	tmplt, _ = template.ParseFiles("index.gohtml")

	var f, _ = os.Create("index.html")

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

	feedsLength := len(feeds)

	var wg sync.WaitGroup
	wg.Add(feedsLength)

	for _, site := range feeds {
		go func(site RssFeed) {
			defer wg.Done()
			fmt.Println("Getting blogs for " + site.NAME)
			feed, _ := gofeed.NewParser().ParseURL(site.RSS)
			if feed == nil {
				fmt.Println("Couldn't fetch blogs for " + site.NAME)
				return
			}
			for _, item := range feed.Items {
				if item.PublishedParsed == nil {
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
			Date: key,
			Blogs: val,
		}
		daten = append(daten, singleDateBlog)
	}
	sort.Slice(daten, func(i, j int) bool {
		return daten[i].Date > daten[j].Date
	})
	generateHtmlFile(daten)
}