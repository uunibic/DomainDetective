package registrars

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
)

type NameCom struct {
    Client   *http.Client
    Username string
    Token    string
}

func NewNameCom(client *http.Client, username, token string) *NameCom {
    return &NameCom{
        Client:   client,
        Username: username,
        Token:    token,
    }
}

type NameComResponse struct {
    Results []struct {
        DomainName     string  `json:"domainName"`
        Purchasable    bool    `json:"purchasable"`
        PurchasePrice  float64 `json:"purchasePrice"`
        RenewalPrice   float64 `json:"renewalPrice"`
    } `json:"results"`
}

func (nc *NameCom) CheckDomainAvailability(domain string) (bool, float64, error) {
    requestURL := "https://api.name.com/v4/domains:checkAvailability"
    requestBody, _ := json.Marshal(map[string][]string{"domainNames": {domain}})

    req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(requestBody))
    if err != nil {
        return false, 0, err
    }

    req.SetBasicAuth(nc.Username, nc.Token)
    req.Header.Add("Content-Type", "application/json")

    resp, err := nc.Client.Do(req)
    if err != nil {
        return false, 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, 0, err
    }

    var response NameComResponse
    err = json.Unmarshal(body, &response)
    if err != nil {
        return false, 0, err
    }

    if len(response.Results) > 0 && response.Results[0].Purchasable {
        return true, response.Results[0].PurchasePrice, nil
    }
    return false, 0, nil
}