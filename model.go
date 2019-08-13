package main

import (
	"net/http"
	"sort"
)

type LookUpResults struct {
	Id      string         `json:"id"`
	Results []LookUpResult `json:"results"`
}

func (*LookUpResults) Render(w http.ResponseWriter, r *http.Request) error {
	// nothing to do
	return nil
}

func (r *LookUpResults) SortResults() {
	sort.Slice(r.Results, func(i, j int) bool {
		return r.Results[i].Number < r.Results[j].Number
	})
}

type LookUpResult struct {
	Number int              `json:"number"`
	IPs    []LookUpResultIp `json:"ips"`
}

type LookUpResultIp struct {
	Address string          `json:"address"`
	Info    *IpLookUpResult `json:"info"`
}
