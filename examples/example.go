package main

import (
	"fmt"
	"os"

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
		username: "API USERNAME",
		apiToken: "API TOKEN",
	}
	client := namecheap.NewClient(account.username, account.apiToken, account.username)

	fmt.Println(Magenta("Namecheap API Example"), Gray(":"), Blue("Requesting Domains"))
	fmt.Println(Gray("=========================================="))
	fmt.Println(Blue("[Namecheap API]"), Gray("Requesting all domains registered by the user"), Green(account.username))

	domains, err := client.DomainsGetList(1, 100)
	fmt.Println(Blue("Requesting first"), Green("100"), Blue("domains from Namecheap API"))
	if err != nil {
		fmt.Println("[Fatal Error]", err)
		os.Exit(1)
	}
	for _, domain := range domains {
		fmt.Println(Blue("" + domain.Name))
	}

	domainCount, err := client.DomainCount()
	if err != nil {
		fmt.Println(Red("[Error]"), ":", err)
	}
	fmt.Println(Blue("[Namecheap API]"), Gray("DomainCount: "), Green(domainCount))

	domains, err := client.DomainsGetCompleteList()
	fmt.Println(Green("Total items received in response: "), len(domains))
	fmt.Println(Blue("Requesting"), Green("ALL"), Blue("domains from Namecheap API"))
	if err != nil {
		fmt.Println("[Fatal Error]", err)
		os.Exit(1)
	}
	for _, domain := range domains {
		fmt.Println(Blue("" + domain.Name))
	}
	fmt.Println(Green("Total items received in response: "), len(domains))
}
