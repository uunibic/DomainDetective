## About

It all began while I was searching for a suitable domain for my portfolio `maheshjandwani.tech` and found it quite expensive (poor me). Luckily, at the same time, I was learning GoLang. This sparked an idea: why not tackle this challenge programmatically? It then emerged as a solution to both a personal need and a programming challenge. So, I wrote a code that helps us check domain availability and compare prices across the below-mentioned registrars, saving us time and money in our domain hunt.

1) GoDaddy
2) Hostinger
3) Name.com
4) DreamHost

## Installation

1. Super easy :) [Download](https://go.dev/doc/install) and install the latest version of Go:

   ```
   $ brew install go
   ```
2. Verify that you've installed Go by opening a command prompt and typing the following command:
   
   ```
   $ go version
   ```
3. Clone the Git repository on your device:

   ```
   $ git clone https://github.com/uunibic/DomainDetective.git && cd DomainDetective
   ```
4. Build the package:

   ```
   $ go build
   ```
5. Use the command to synchronize with the actual dependencies used in the codebase:

   ```
   $ go mod init
   ```
6. Set the env variables:

   ```
   export GODADDY_API_KEY="<value-here>"
   export GODADDY_API_SECRET="<value-here>"
   export HOSTINGER_AUTH_TOKEN="www.hostinger.com"
   export NAMECOM_SECRET_KEY="<value-here>"
   export NAMECOM_USER_NAME="<value-here>"
   ```

## Usage

- For querying a single domain: `(-d)`

  ```
  ./DomainDetective -d <domain-name-here>
  ```
- For querying multiple domains using text file: `(-f)`

  ```
  ./DomainDetective -f <txt-file-name-here>
  ```

## Screenshots

1. Single Domain Lookup:

   <img width="692" alt="Screenshot1" src="https://github.com/uunibic/DomainDetective/assets/64989501/b637dd82-a8ac-40d4-85ba-5695539bd563">

3. Multi-Domain Lookup:

   <img width="684" alt="Screenshot2" src="https://github.com/uunibic/DomainDetective/assets/64989501/1bfdd7f1-1131-41bd-8118-7c027448800a">

## Contributions

Although this was a learning project, I would love to invite you guys to join me in making domain hunting a breeze with DomainDetective! Your contribution, big or small, plays a part in this exciting journey. Thank you for your interest and for checking out the tool.
