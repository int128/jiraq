package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/xerrors"

	jira "github.com/andygrunwald/go-jira"
)

type cmdOptions struct {
	url      string
	username string
	password string
	jql      string
	timezone string
}

func run(osArgs []string) error {
	var o cmdOptions
	f := flag.NewFlagSet(osArgs[0], flag.ContinueOnError)
	f.StringVar(&o.url, "url", "", "Jira URL")
	f.StringVar(&o.username, "username", "", "Username of Jira user")
	f.StringVar(&o.password, "password", "", "Password of Jira user")
	f.StringVar(&o.jql, "jql", "", "JQL")
	f.StringVar(&o.timezone, "tz", "Local", "Timezone")
	if err := f.Parse(os.Args[1:]); err != nil {
		return xerrors.Errorf("invalid argument: %s", err)
	}

	tp := jira.BasicAuthTransport{
		Username: o.username,
		Password: o.password,
	}
	c, err := jira.NewClient(tp.Client(), o.url)
	if err != nil {
		return xerrors.Errorf("could not connect to Jira: %w", err)
	}
	log.Printf("Searching issues: %s", o.jql)
	issues, _, err := c.Issue.Search(o.jql, &jira.SearchOptions{
		MaxResults: 100,
		Fields:     []string{"key", "summary"},
		Expand:     "changelog",
	})
	if err != nil {
		return xerrors.Errorf("could not search the issues: %w", err)
	}
	localLocation, err := time.LoadLocation(o.timezone)
	if err != nil {
		return xerrors.Errorf("could not load the location: %s", err)
	}
	for _, issue := range issues {
		fmt.Printf("%s,%s\n", issue.Key, issue.Fields.Summary)
		if issue.Changelog != nil {
			for _, history := range issue.Changelog.Histories {
				created, err := history.CreatedTime()
				if err != nil {
					return xerrors.Errorf("invalid created time of issue %s: %s", issue.Key, err)
				}
				createdLocal := created.In(localLocation)
				for _, item := range history.Items {
					if item.Field == "status" {
						fmt.Printf(",%s,%s\n", createdLocal, item.ToString)
					}
				}
			}
		}
	}
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("error: %s", err)
	}
}
