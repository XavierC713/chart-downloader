package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	// config
	path   string
	source ChartSource

	// aviaplanner specific
	aviaToken string
	aviaPid   string

	// globals
	reader *bufio.Reader
)

func main() {
	reader = bufio.NewReader(os.Stdin)
	userConfig()

	icao := prompt("Enter airport ICAO code: ")
	switch source {
	case SOURCE_FAA:
		charts := getFAA(icao)
		for _, chart := range charts {
			fmt.Printf("\nDownloading %s\n", chart.ChartName)
			downloadChart(chart.PdfPath, fmt.Sprintf("%s - %s.pdf", chart.IcaoIdent, chart.ChartName))
		}
	case SOURCE_AVIAPLANNER:
		getLIDO(icao)
	}
}

// configure all options according to user input, ran at start of program
func userConfig() {
	chartSource := promptWithOptions("Which source should be used for charts?", []string{
		"FAA",
		"AviaPlanner",
	})
	switch chartSource {
	case 1:
		source = SOURCE_FAA
	case 2:
		source = SOURCE_AVIAPLANNER
	}

	if source == SOURCE_AVIAPLANNER {
		// need to get token & pid
		fmt.Println("Due to limiations with AviaPlanner, certain information is needed to retrieve charts from the platform. This includes two cookies called the token and pid which help AviaPlanner authenticate you and verify your subscription status. These cookies should NOT be shared publicly as they could allow someone to use your account to make requests to the AviaPlanner platform. This software only uses these cookies to make requests and does not store them, this can be verified by looking throught the source code. If you do not know how to obtain these cookies, please refer to this guide: https://github.com/XavierC713/chart-downloader/blob/master/aviaplanner_cookie_guide.md")
		aviaToken = prompt("Token: ")
		aviaPid = prompt("PID: ")
	}

	path = prompt("Where should charts be downloaded? ")
	if strings.HasPrefix(path, "~") {
		path = os.ExpandEnv(strings.Replace(path, "~", "$HOME", 1))
	}
}

func downloadChart(url string, fileName string) {
	out, err := os.Create(fmt.Sprintf("%s/%s", path, fileName))
	if err != nil {
		fmt.Println(err)
		panic("Could not create file in provided path, aborting...")
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("\nFailed to fetch %s", url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("\nError when fetching %s, status: %s", url, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to write file %s, aborting...", fileName))
	}
}
