# go-namecheap

A Go library for using [the Namecheap XML REST API](https://www.namecheap.com/support/api/intro.aspx).

## Examples

```go
package main

import (
	"fmt"

	. "github.com/hackwave/color"
	namecheap "github.com/scrambleshell/namecheap-go"
)

type Domain struct {
	name string
}

type Account struct {
	username string
	// Use Memguard to protect the password by clearing the memory after using it
	password string
	apiToken string
	domains  []Domain
}

func main() {
	// TODO: Use surf (and otto if necesssary) to login and generate an API token if neccessary and 

	account := Account{
		username: "kosmosblack",
		apiToken: "{API_TOKEN}",
	}

	client := namecheap.NewClient(account.username, account.apiToken, account.username)

	fmt.Println(Magenta("Namecheap Domains"))
	fmt.Println(Gray("================="))

	fmt.Println(Gray("Looking up domains for the API account: "), Green(account.username))
	// Get a list of your domains
	domains, err := client.DomainsGetList(1, 100)
	if err != nil {
		fmt.Println(Red("[Error]"), err)
	}

	fmt.Println(Gray("Number of found domains: "), Green(len(domains)))
	for _, domain := range domains {
		fmt.Printf("Domain: %+v\n\n", domain.Name)
	}
}
```

