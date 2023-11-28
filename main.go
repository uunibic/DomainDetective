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
    dreamhost := registrars.NewDreamHost(client, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
    namecom := registrars.NewNameCom(client, os.Getenv("NAMECOM_USER_NAME"), os.Getenv("NAMECOM_SECRET_KEY"))

    if *domain != "" {
        processDomain(*domain, godaddy, hostinger, dreamhost, namecom)
    }

    if *filename != "" {
        file, err := os.Open(*filename)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            processDomain(scanner.Text(), godaddy, hostinger, dreamhost, namecom)
        }

        if err := scanner.Err(); err != nil {
            log.Fatal(err)
        }
    }
}

func processDomain(domainName string, godaddy, hostinger, dreamhost, namecom registrars.Registrar) {

    availableGoDaddy, priceGoDaddy, err := godaddy.CheckDomainAvailability(domainName)
    if err != nil {
        log.Printf("Error checking domain %s on GoDaddy: %v\n", domainName, err)
        return
    } else if !availableGoDaddy {
        fmt.Printf("[Registrar] %s - Domain not available ❌\n", domainName)
        return
    } else {
        fmt.Printf("[GoDaddy] %s - Available for $%.2f (%d year) ✅\n", domainName, priceGoDaddy, 1)
    }

    results := make(chan string)
    go checkAndPrintAvailabilityConcurrent(hostinger, "Hostinger", domainName, results)
    go checkAndPrintAvailabilityConcurrent(dreamhost, "DreamHost", domainName, results)
    go checkAndPrintAvailabilityConcurrent(namecom, "Name.com", domainName, results)

    for i := 0; i < 3; i++ {
        result := <-results
        if result != "" {
            fmt.Print(result)
        }
    }
    close(results)
}

func checkAndPrintAvailabilityConcurrent(registrar registrars.Registrar, registrarName, domainName string, results chan<- string) {
    available, price, err := registrar.CheckDomainAvailability(domainName)
    if err != nil {
        log.Printf("Error checking domain %s on %s: %v\n", domainName, registrarName, err)
        results <- ""
    } else if available {
        results <- fmt.Sprintf("[%s] %s - Available for $%.2f (%d year) ✅\n", registrarName, domainName, price, 1)
    } else {
        results <- ""
    }
}