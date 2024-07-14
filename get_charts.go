package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type FAAChart struct {
	State       string `json:"state"`
	StateFull   string `json:"state_full"`
	City        string `json:"city"`
	AirportName string `json:"airport_name"`
	Military    string `json:"military"`
	FaaIdent    string `json:"faa_ident"`
	IcaoIdent   string `json:"icao_ident"`
	ChartSeq    string `json:"chart_seq"`
	ChartCode   string `json:"chart_code"`
	ChartName   string `json:"chart_name"`
	PdfName     string `json:"pdf_name"`
	PdfPath     string `json:"pdf_path"`
}

type FAAChartResponse map[string][]FAAChart

type AviaPlannerAirportResponse struct {
	Res    int      `json:"res"`
	Err    int      `json:"err"`
	Coords []string `json:"coords"`
	HTML   string   `json:"html"`
}

// fetch all FAA charts for an airport from aviationapi
func getFAA(icao string) []FAAChart {
	resp, err := http.Get(fmt.Sprintf("https://api.aviationapi.com/v1/charts?apt=%s", icao))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result FAAChartResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	return result[strings.ToUpper(icao)]
}

// fetch all LIDO charts for an airport from AviaPlanner
// IMPORTANT: aviaplanner does not have a public api so this method uses internal endpoints and scrapes the returned data, these endpoints may change in the future causing this to not work as expected
func getLIDO(icao string) map[string]string {
	data := url.Values{}
	data.Set("icao", icao)

	req, err := http.NewRequest(http.MethodPost, "https://web.aviaplanner.com/ajax/?type=port-info", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: aviaToken,
		Path:  "/",
	})
	req.AddCookie(&http.Cookie{
		Name:  "pid",
		Value: aviaPid,
		Path:  "/",
	})
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0") // impersonate firefox because aviaplanner will send a 403 otherwise
	req.Header.Add("Priority", "u=0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result AviaPlannerAirportResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	regex := *regexp.MustCompile(`<div class="procedureLine"> <div class="dataProcedure wa"><div class="wp100"><span class="info">.*?</span><span class="charts">(.*?)</span></div></div> <div class="operationButtons wa"> <a href="javascript:Planner\.showChart\((\d*?)\)" class="iBut" title="View Lido chart"><i class="apb chart"></i></a> </div> </div>`)
	matches := regex.FindAllStringSubmatch(result.HTML, -1)
	charts := map[string]string{}
	for _, match := range matches {
		chartName := match[1]
		chartId := match[2]
		charts[chartName] = chartId
	}
	return charts
}
