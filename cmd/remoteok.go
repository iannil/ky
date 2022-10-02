package cmd

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

var cmdRemoteOK = &cobra.Command{
	Use:   "remoteok",
	Aliases: []string{"ro"},
	Short: "从remoteok.io获取最新的远程工作",
	Run: func(cmd *cobra.Command, args []string) {
		start := "https://remoteok.io/worldwide"
		if len(args) == 1 {
			start = "https://remoteok.io/remote-"+args[0]+"-jobs?location=worldwide"
		}

		table := pterm.TableData{
			{"Name", "Tags", "Time", "Link"},
		}
		c := colly.NewCollector(colly.CacheDir("./cache"))

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("tr.job", func(e *colly.HTMLElement) {
			var name, fullname, link, time, tags, isClosed string
			e.ForEach("td", func(i int, element *colly.HTMLElement) {
				if i == 1 {
					name = element.ChildText("span.companyLink")
					fullname = element.ChildText("a.preventLink")
					link = element.ChildAttr("a.preventLink", "href")
					isClosed = element.ChildText("span.closed")
				}
				if i == 3 {
					tagsList := element.ChildTexts("a.action-add-tag")
					tags = strings.Join(tagsList, ",")
				}
				if i == 4 {
					re := regexp.MustCompile("([0-9]+).?")
					time = re.FindString(element.ChildText("a"))
				}
			})
			if isClosed != "closed" {
				table = append(table, []string{name+","+fullname, tags, time, link})
			}
		})

		c.Visit(start)

		defer func() {
			pterm.DefaultTable.WithHasHeader().WithData(table).Render()
		}()
	},
}
