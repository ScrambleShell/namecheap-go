# go-namecheap

A Go library for using [the Namecheap XML REST API](https://www.namecheap.com/support/api/intro.aspx).

## Examples

```go
package main

import (
	"fmt"
	"os"

	namecheap "github.com/scrambleshell/namecheap-go"
)

type Account struct {
	username string
	apiToken string
}

func main() {
	fmt.Println("Namecheap Domains")
	fmt.Println("=================")

	account := Account{
		username: "{API_USERNAME}",
		apiToken: "{API_TOKEN}",
	}
	client := namecheap.NewClient(account.username, account.apiToken, account.username)

	fmt.Println("[Namecheap API] Requesting all domains registered by the user", account.username)
	domains, err := client.DomainsGetList(1, 100)
	if err != nil {
		fmt.Println("[Fatal Error]", err)
		os.Exit(1)
	}

	for _, domain := range domains {
		fmt.Printf("Domain: %+v\n\n", domain.Name)
	}
}
```

