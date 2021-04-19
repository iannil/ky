package cmd

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
	"strings"
	"sync/atomic"
)

var cmdNyaa = &cobra.Command{
	Use:   "nyaa",
	Short: "根据番号从sukebei.nyaa.si获取磁力链接",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fan := args[0]

		var ops uint64

		c := colly.NewCollector(colly.CacheDir("./cache"))

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if strings.Contains(href, "magnet") {
				atomic.AddUint64(&ops, 1)
				fmt.Printf("[%d]%s\n", ops, href)
			}
		})

		c.Visit("https://sukebei.nyaa.si/?f=0&c=0_0&q=" + fan + "&s=downloads&o=desc")
	},
}
