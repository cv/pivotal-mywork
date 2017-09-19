package main

import (
	"fmt"
	"os"

	"time"

	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/cv/go-pivotaltracker/v5/pivotal"
)

// 35 days ago
const SINCE = -35 * 24 * time.Hour

func main() {
	token := os.Getenv("PIVOTAL_TOKEN")
	if token == "" {
		fmt.Println("PIVOTAL_TOKEN environment variable not set")
		os.Exit(-1)
	}

	client := pivotal.NewClient(token)
	me, _, err := client.Me.Get()
	if err != nil {
		return
	}

	since := time.Now().Add(SINCE).Format("01/02/2006")
	fmt.Printf("# Listing stories for %s, aka %s\n",
		color.CyanString(me.Name),
		color.CyanString(me.Initials),
	)

	projects, _, err := client.Projects.List()
	if err != nil {
		fmt.Printf("Error fetching projects: %v", err)
		os.Exit(-2)
	}

	allStories := []*pivotal.Story{}
	for _, project := range projects {
		stories, err := client.Stories.List(project.Id, fmt.Sprintf("owner:%d includedone:true updated_since:%s", me.Id, since))
		if err != nil {
			fmt.Printf("Error fetching stories for project %s: %v", project.Id, err)
			os.Exit(-2)
		}

		allStories = append(allStories, stories...)
	}

	sort.Slice(allStories, func(i, j int) bool {
		return allStories[i].UpdatedAt.After(*allStories[j].UpdatedAt)
	})

	byAcceptedAt := map[string][]*pivotal.Story{}

	fmt.Println()
	fmt.Printf("# Stories in the 'My Work' queue updated since %s:\n", color.HiCyanString(since))

	for _, story := range allStories {
		updated := "-"
		if story.UpdatedAt != nil {
			updated = story.UpdatedAt.Format("02/Jan/2006")
		}

		done := "-"
		if story.AcceptedAt != nil {
			done = story.AcceptedAt.Format("02/Jan/2006")
			byAcceptedAt[done] = append(byAcceptedAt[done], story)
		} else {
			fmt.Printf("%s %s %s %s %s %s\n",
				color.HiBlueString(fmt.Sprint(story.Id)),
				color.YellowString(story.Type),
				color.MagentaString(story.State),
				color.GreenString(updated),
				color.HiGreenString(done),
				story.Name,
			)
		}
	}

	fmt.Println()
	fmt.Printf("# Stories completed since %s:\n", color.HiCyanString(since))

	for k, v := range byAcceptedAt {
		fmt.Printf("%s: ", color.HiGreenString(k))
		for _, s := range v {
			fmt.Printf("%s %s %s: %s\n             ",
				color.HiBlueString(fmt.Sprint(s.Id)),
				color.YellowString(s.Type),
				color.MagentaString(s.State),
				strings.TrimSpace(s.Name))
		}
		fmt.Println()
	}
}
