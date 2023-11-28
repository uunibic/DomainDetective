package registrars

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
)

type DreamHost struct {
    Client    *http.Client
    UserAgent string
}

func NewDreamHost(client *http.Client, userAgent string) *DreamHost {
    return &DreamHost{
        Client:    client,
        UserAgent: userAgent,
    }
}

type DreamHostResponse struct {
    Available     string  `json:"available"`
    IsPremium     string  `json:"is_premium"`
    PremiumPrice  *float64 `json:"premium_price"`
    Price         float64 `json:"price"`
    RenewPrice    float64 `json:"renew_price"`
    TLD           string  `json:"tld"`
}

func (dh *DreamHost) CheckDomainAvailability(domain string) (bool, float64, error) {
    url := fmt.Sprintf("https://marketing-api-aws.dreamhost.io/ajax.cgi?cmd=domreg-availability&domain=%s&pricing_wanted=1&callback=jsonp", domain)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false, 0, err
    }

    req.Header.Set("User-Agent", dh.UserAgent)

    resp, err := dh.Client.Do(req)
    if err != nil {
        return false, 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, 0, err
    }

    trimmedResponse := strings.TrimSuffix(strings.TrimPrefix(string(body), "jsonp("), ")")

    var response DreamHostResponse
    err = json.Unmarshal([]byte(trimmedResponse), &response)
    if err != nil {
        return false, 0, err
    }

    available := response.Available == "true"
    return available, response.Price, nil
}