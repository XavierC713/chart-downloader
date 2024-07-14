package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func getFAA(icao string) []FAAChart {
	resp, err := http.Get(fmt.Sprintf("https://api.aviationapi.com/v1/charts?apt=%s", icao))
	if err != nil {
		panic("No response from API, aborting...")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result FAAChartResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic("Unexpected format returned by API, aborting...")
	}
	return result[strings.ToUpper(icao)]
}
