package cmd

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
	"ky/internal"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var cmdWnacg = &cobra.Command{
	Use:   "wnacg [album id] [path to save]",
	Short: "下载wnacg.org的漫画",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := args[0]

		output := "./dist"
		title := fmt.Sprintf("%d", time.Now().Unix())

		if len(args) > 1 {
			output = args[1]
		}

		output = fmt.Sprintf("%s/%s/", output, title)
		err := os.MkdirAll(output, 0755)
		if err != nil {
			panic(err)
		}

		var images []string
		cover := "this-is-cover"

		var ops uint64

		c := colly.NewCollector(colly.CacheDir("./cache"))

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("h2", func(e *colly.HTMLElement) {
			h2 := e.Text
			if h2 != "" {
				title = h2
			}
		})

		c.OnHTML("img", func(e *colly.HTMLElement) {
			src := e.Attr("src")

			if strings.HasPrefix(src, "////t") {
				segs := strings.Split(src, "/")
				cover = segs[len(segs)-1]
				fmt.Println("Cover", cover)
			}
		})

		c.OnHTML("img.photo", func(e *colly.HTMLElement) {
			src := e.Attr("src")

			fmt.Println("Found: ", src)
			if strings.Contains(src, "/data/") && !strings.Contains(src, cover) {
				src = "https:" + src

				segs := strings.Split(src, "/")
				filename := segs[len(segs)-1]

				dest := output + "/" + fmt.Sprintf("%d-", ops) + filename
				images = append(images, dest)
				err := internal.DownloadFile(src, dest)
				if err != nil {
					fmt.Println("error: ", err)
				}
				atomic.AddUint64(&ops, 1)
			}
		})

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if strings.Contains(href, "/photos-index-page") && !strings.Contains(href, "photos-index-page-1-") {
				e.Request.Visit(href)
			}
			if strings.Contains(href, "/photos-view-id-") {
				e.Request.Visit(href)
			}
		})

		c.Visit("https://wnacg.org/photos-index-aid-" + start + ".html")

		defer func() {
			fmt.Println("Zipping...")
			archiver.Archive(images, fmt.Sprintf("%s%s.zip", output, title))
		}()
	},
}
