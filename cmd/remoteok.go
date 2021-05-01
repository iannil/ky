package cmd

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"strings"
)

var cmdRemoteOK = &cobra.Command{
	Use:   "remoteok",
	Aliases: []string{"ro"},
	Short: "从remoteok.io获取最新的PHP远程工作",
	Run: func(cmd *cobra.Command, args []string) {
		table := pterm.TableData{
			{"Name", "Tags", "Time", "Link"},
		}
		c := colly.NewCollector(colly.CacheDir("./cache"))

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("tr.job.remoteok-original", func(e *colly.HTMLElement) {
			var name, fullname, link, time, tags string
			e.ForEach("td", func(i int, element *colly.HTMLElement) {
				if i == 1 {
					name = element.ChildText("span.companyLink")
					fullname = element.ChildText("a.preventLink")
					link = element.ChildAttr("a.preventLink", "href")
				}
				if i == 3 {
					tagsList := element.ChildTexts("a.action-add-tag")
					tags = strings.Join(tagsList, ",")
				}
				if i == 4 {
					time = element.ChildText("a")
				}
			})
			table = append(table, []string{name+","+fullname, tags, time, link})
		})

		c.Visit("https://remoteok.io/remote-php-jobs?location=worldwide")

		defer func() {
			pterm.DefaultTable.WithHasHeader().WithData(table).Render()
		}()
	},
}
