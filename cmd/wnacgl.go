package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
)

var cmdWnacgl = &cobra.Command{
	Use:   "wnacgl [start]",
	Short: "下载wnacg.org的漫画",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := args[0]

		c := colly.NewCollector(colly.CacheDir("./cache"))

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML(".info", func(e *colly.HTMLElement) {
			title := strings.Join(strings.Fields(e.DOM.Find(".title").Text()), "")
			href, _ := e.DOM.Find(".title a").Attr("href")
			link := strings.Join(strings.Fields(href), "")
			info := strings.Join(strings.Fields(e.DOM.Find(".info_col").Text()), "")

			r, _ := regexp.Compile("([0-9]+)張照片")
			res := r.FindStringSubmatch(info)
			infocount, _ := strconv.Atoi(res[1])
			if infocount >= 100 || strings.Contains(strings.ToLower(title), "vol") {
				fmt.Printf("%s,%s,%s,%d\n", title, info, link, infocount)
			}
		})

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if strings.Contains(href, "/albums-index-page") && !strings.Contains(href, "/albums-index-page-1-") {
				e.Request.Visit(href)
			}
		})

		c.Visit("https://wnacg.org/albums-index-cate-" + start + ".html")
	},
}
