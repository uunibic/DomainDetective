package registrars

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
)

type Hostinger struct {
    Client    *http.Client
    AuthToken string
}

func NewHostinger(client *http.Client, authToken string) *Hostinger {
    return &Hostinger{
        Client:    client,
        AuthToken: authToken,
    }
}

func (h *Hostinger) CheckDomainAvailability(domain string) (bool, float64, error) {
    parts := strings.SplitN(domain, ".", 2)
    if len(parts) != 2 {
        return false, 0, fmt.Errorf("invalid domain format: %s", domain)
    }

    domainReq := HostingerDomainRequest{
        DomainName: parts[0],
        TLD:        parts[1],
    }
    requestBody, err := json.Marshal(domainReq)
    if err != nil {
        return false, 0, err
    }

    req, err := http.NewRequest("POST", "https://websites-api.hostinger.com/api/domain/single-domain-search", bytes.NewBuffer(requestBody))
    if err != nil {
        return false, 0, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.AuthToken))

    resp, err := h.Client.Do(req)
    if err != nil {
        return false, 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, 0, err
    }

    var domainInfo HostingerDomainResponse
    err = json.Unmarshal(body, &domainInfo)
    if err != nil {
        return false, 0, err
    }

    return domainInfo.Data.Result.Available, domainInfo.Data.Result.Product.Price.Purchase, nil
}

type HostingerDomainRequest struct {
    DomainName string `json:"domain_name"`
    TLD        string `json:"tld"`
}

type HostingerDomainResponse struct {
    Data struct {
        Result struct {
            Available  bool   `json:"available"`
            DomainName string `json:"domain_name"`
            Product    struct {
                Price struct {
                    Purchase float64 `json:"purchase"`
                } `json:"price"`
            } `json:"product"`
        } `json:"result"`
    } `json:"data"`
}