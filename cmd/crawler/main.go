package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-labs/internal/crawler"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var source, destinationDir string

func init() {
	rootCmd.Flags().StringVarP(&source, "source", "s", "", "URL to crawl.")
	rootCmd.Flags().StringVarP(
		&destinationDir,
		"destination",
		"d",
		crawler.DestinationDir,
		"Directory to store downloaded contents.",
	)

	_ = rootCmd.MarkFlagRequired(source)
}

var rootCmd = &cobra.Command{
	Use:   "crawler",
	Short: "Crawler is a recursive, mirroring web crawler.",
	Long:  `A simple web crawl that downloads and crawls all relative links for the provided URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			client   = &http.Client{}
			signalCh = make(chan os.Signal)
		)

		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-signalCh
			_, _ = fmt.Fprintf(os.Stderr, "Signal interrupt - exiting\n")
			os.Exit(1)
		}()

		crawl, err := crawler.NewCrawler(destinationDir, client)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		var wg sync.WaitGroup
		wg.Add(1)

		crawl.Crawl(source, 4, crawl, &wg)
		wg.Wait()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "There was an error while executing your command '%s'\n", err)
		os.Exit(1)
	}
}
