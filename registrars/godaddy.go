package registrars

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type GoDaddy struct {
    Client     *http.Client
    APIKey     string
    APISecret  string
}

type DomainInfo struct {
    Available  bool   `json:"available"`
    Currency   string `json:"currency"`
    Definitive bool   `json:"definitive"`
    Domain     string `json:"domain"`
    Period     int    `json:"period"`
    Price      int    `json:"price"`
}

func NewGoDaddy(client *http.Client, apiKey, apiSecret string) *GoDaddy {
    return &GoDaddy{
        Client:    client,
        APIKey:    apiKey,
        APISecret: apiSecret,
    }
}

func (g *GoDaddy) CheckDomainAvailability(domain string) (bool, float64, error) {
    url := fmt.Sprintf("https://api.godaddy.com/v1/domains/available?domain=%s&checkType=FAST&forTransfer=false", domain)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false, 0, err
    }

    authHeader := fmt.Sprintf("sso-key %s:%s", g.APIKey, g.APISecret)
    req.Header.Set("accept", "application/json")
    req.Header.Set("Authorization", authHeader)

    resp, err := g.Client.Do(req)
    if err != nil {
        return false, 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, 0, err
    }

    var domainInfo DomainInfo
    err = json.Unmarshal(body, &domainInfo)
    if err != nil {
        return false, 0, err
    }

    priceInUSD := float64(domainInfo.Price) / 1000000
    return domainInfo.Available, priceInUSD, nil
}