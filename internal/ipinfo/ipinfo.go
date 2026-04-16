// Package ipinfo provides lightweight metadata lookup for IP addresses,
// returning geolocation and organisation details for alert enrichment.
package ipinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Info holds metadata returned for an IP address.
type Info struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Org     string `json:"org"`
}

// String returns a human-readable summary.
func (i Info) String() string {
	if i.City == "" && i.Country == "" {
		return i.IP
	}
	return fmt.Sprintf("%s (%s, %s) %s", i.IP, i.City, i.Country, i.Org)
}

// Lookup fetches metadata for the given IP using the ipinfo.io public API.
// An empty token is accepted; rate limits will apply.
type Lookup struct {
	client  *http.Client
	baseURL string
	token   string
}

// New returns a Lookup with the given API token (may be empty).
func New(token string) *Lookup {
	return &Lookup{
		client:  &http.Client{Timeout: 5 * time.Second},
		baseURL: "https://ipinfo.io",
		token:   token,
	}
}

// Get returns Info for the given IP address.
func (l *Lookup) Get(ip string) (Info, error) {
	url := fmt.Sprintf("%s/%s/json", l.baseURL, ip)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Info{}, err
	}
	if l.token != "" {
		req.Header.Set("Authorization", "Bearer "+l.token)
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return Info{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf("ipinfo: unexpected status %d", resp.StatusCode)
	}
	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return Info{}, err
	}
	return info, nil
}
