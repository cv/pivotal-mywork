package cmd

import (
	"fmt"
	"os"

	"sort"

	"strings"

	"github.com/cv/go-pivotaltracker/v5/pivotal"
	"github.com/jaytaylor/html2text"
	"github.com/olekukonko/tablewriter"
	"github.com/russross/blackfriday"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Print information about a card",
	Run: func(_ *cobra.Command, args []string) {
		token := os.Getenv("PIVOTAL_TOKEN")
		if token == "" {
			fmt.Println("PIVOTAL_TOKEN environment variable not set")
			os.Exit(-1)
		}

		client := pivotal.NewClient(token)

		projects, _, err := client.Projects.List()
		if err != nil {
			fmt.Printf("Error fetching projects: %v", err)
			os.Exit(-2)
		}

		if len(args) != 1 {
			fmt.Println("Usage: pivotal-get <STORY ID>")
			os.Exit(-3)
		}

		id := args[0]
		if id == "" {
			fmt.Println("No story ID given")
			os.Exit(-3)
		}

		allStories := []*pivotal.Story{}
		for _, project := range projects {
			stories, err := client.Stories.List(project.Id, fmt.Sprintf("includedone:true id:\"%s\"", id))
			if err != nil {
				fmt.Printf("Error fetching stories for project %s: %v", project.Id, err)
				os.Exit(-2)
			}

			allStories = append(allStories, stories...)
		}

		sort.Slice(allStories, func(i, j int) bool {
			return allStories[i].UpdatedAt.After(*allStories[j].UpdatedAt)
		})

		for _, story := range allStories {

			table := tablewriter.NewWriter(os.Stdout)

			table.SetColWidth(80)
			table.SetAutoWrapText(true)
			table.SetAlignment(tablewriter.ALIGN_LEFT)

			var project *pivotal.Project
			for _, project = range projects {
				if project.Id == story.ProjectId {
					break
				}
			}

			table.Append([]string{"Story", fmt.Sprintf("%d (%s, %s) in %s", story.Id, story.Type, story.State, project.Name)})

			if story.Estimate == nil {
				table.Append([]string{"Estimate", "N/A"})
			} else {
				table.Append([]string{"Estimate", fmt.Sprintf("%.0f", *story.Estimate)})
			}

			if story.CreatedAt != nil {
				table.Append([]string{"Created", story.CreatedAt.Format("02/Jan/2006")})
			}
			if story.UpdatedAt != nil {
				table.Append([]string{"Updated", story.UpdatedAt.Format("02/Jan/2006")})
			}
			if story.AcceptedAt != nil {
				table.Append([]string{"Accepted", story.AcceptedAt.Format("02/Jan/2006")})
			}
			if story.Deadline != nil {
				table.Append([]string{"Deadline", story.Deadline.Format("02/Jan/2006")})
			}

			table.Append([]string{"", ""})
			table.Append([]string{"Title", story.Name})
			table.Append([]string{"URL", story.URL})
			table.Append([]string{"", ""})

			table.Append([]string{"Description", renderClean(story.Description)})

			table.Render()
		}
	},
}

func renderClean(text string) string {
	txt, _ := html2text.FromString(string(blackfriday.MarkdownBasic([]byte(text))))
	return strings.Replace(strings.Replace(txt, " )", ")", -1), "( ", "(", -1)
}
