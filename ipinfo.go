package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const baseIpInfoUrl = "https://ipinfo.io"

type IpInfoClient struct {
	Key string
}

func NewIpInfoClient(key string) *IpInfoClient {
	return &IpInfoClient{Key: key}
}

type IpLookUpResult struct {
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Postal   string `json:"postal"`
	ASN      struct {
		ASN    string `json:"asn"`
		Name   string `json:"name"`
		Domain string `json:"domain"`
		Route  string `json:"route"`
		Type   string `json:"type"`
	} `json:"asn"`
	Bogon bool `json:"bogon"`
}

func (ip *IpInfoClient) LookUpIp(addr string) (*IpLookUpResult, error) {
	resp, err := http.Get(baseIpInfoUrl + "/" + addr)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ipinfo")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid http status response %d", resp.StatusCode)
	}

	var result IpLookUpResult
	return &result, errors.Wrap(json.NewDecoder(resp.Body).Decode(&result), "unable to decode json response")
}
