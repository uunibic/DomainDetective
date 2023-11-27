package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "bufio"
    "github.com/uunibic/DomainDetective/registrars"
)

func main() {
    domain := flag.String("d", "", "Specify the domain to check")
    filename := flag.String("f", "", "Specify the file with domain list")
    flag.Parse()

    client := &http.Client{}

    godaddy := registrars.NewGoDaddy(client, os.Getenv("GODADDY_API_KEY"), os.Getenv("GODADDY_API_SECRET"))
    hostinger := registrars.NewHostinger(client, os.Getenv("HOSTINGER_AUTH_TOKEN"))

    if *domain != "" {
        processDomain(*domain, godaddy, hostinger)
    }

    if *filename != "" {
        file, err := os.Open(*filename)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            processDomain(scanner.Text(), godaddy, hostinger)
        }

        if err := scanner.Err(); err != nil {
            log.Fatal(err)
        }
    }
}

func processDomain(domainName string, godaddy, hostinger registrars.Registrar) {
    availableGoDaddy, priceGoDaddy, err := godaddy.CheckDomainAvailability(domainName)
    if err != nil {
        log.Printf("Error checking domain %s on GoDaddy: %v\n", domainName, err)
        return
    }

    if !availableGoDaddy {
        fmt.Printf("[Registrar] %s - Domain not available ❌\n", domainName)
        return
    }

    fmt.Printf("[GoDaddy] %s - Available for $%.2f (%d year) ✅\n", domainName, priceGoDaddy, 1)

    availableHostinger, priceHostinger, err := hostinger.CheckDomainAvailability(domainName)
    if err != nil {
        log.Printf("Error checking domain %s on Hostinger: %v\n", domainName, err)
        return
    }

    if availableHostinger {
        fmt.Printf("[Hostinger] %s - Available for $%.2f (%d year) ✅\n", domainName, priceHostinger, 1)
    }
}