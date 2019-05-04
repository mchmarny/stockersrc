package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	twitterSourceFormat = "https://twitter.com/%s/status/%s"
	twitterIDFormat     = "tw-%s"
)

var (
	consumerKey    = mustEnvVar("T_CONSUMER_KEY", "")
	consumerSecret = mustEnvVar("T_CONSUMER_SECRET", "")
	accessToken    = mustEnvVar("T_ACCESS_TOKEN", "")
	accessSecret   = mustEnvVar("T_ACCESS_SECRET", "")

	// validation expressions
	spaceReg = regexp.MustCompile(`\s+`)
	spaceCr  = regexp.MustCompile(`^[\r\n]+|\.|[\r\n]+$`)
)

// provide runs twitter search and publishes to provided sinker
func provide(ctx context.Context, co *Company, out chan<- *TextContent) {

	// twitter client config
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	twClient := twitter.NewClient(httpClient)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(t *twitter.Tweet) {
		logger.Printf("Got tweet: %s for %s", t.IDStr, co.Symbol)

		createdTime, err := t.CreatedAtTime()
		if err != nil {
			logger.Printf("Error while parsing created at: %v", err)
			createdTime = time.Now()
		}

		tc := &TextContent{
			Symbol:    co.Symbol,
			ID:        fmt.Sprintf(twitterIDFormat, t.IDStr),
			CreatedAt: createdTime,
			Author:    strings.ToLower(t.User.ScreenName),
			Lang:      t.Lang,
			Source:    fmt.Sprintf(twitterSourceFormat, t.User.ScreenName, t.IDStr),
			Content:   clean(t.Text),
		}

		// TODO" Externalize lang support
		if tc.Lang == "en" {
			out <- tc
		}
	}

	query := []string{co.Symbol}
	alias := strings.Split(co.Aliases, ",")
	query = append(query, alias...)
	logger.Printf("Query: %s", query)

	params := &twitter.StreamFilterParams{
		Track:         query,
		FilterLevel:   "none",
		StallWarnings: twitter.Bool(true),
		Language:      []string{"en"},
	}

	stream, err := twClient.Streams.Filter(params)
	if err != nil {
		logger.Printf("Error on filter create: %v", err)
		return
	}

	logger.Printf("Starting tweet streamming for: %+v", co)
	demux.HandleChan(stream.Messages)

}

func clean(txt string) string {
	txt = spaceReg.ReplaceAllString(txt, " ")
	txt = strings.Trim(txt, " ")
	txt = spaceCr.ReplaceAllString(txt, " ")
	return txt
}
