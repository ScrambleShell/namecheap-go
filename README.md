# Namecheap Go XML API Client
A Go client library that aims to be consistent, provide complete functionality (it is not currently complete), and additionally provide a thin abstraction layer for improving ease-of-use and access to the Namecheap XML [API](https://www.namecheap.com/support/api/intro.aspx).

A minimalistic Namecheap API that is not currently yet feature complete and there is not yet documentation detailing feature completeness or a roadmap outlining a path to feature completeness. The library in additional to providing a complete implementation of the XML API, it will provide a thin abstraction layer to simplify use, abstract additional features to make using the API more intuitive. As the library evolves, the function names will be altered to make use more intuitive and code using the library more readable but the original functions names will still be accessible and usable by aliasing the altered names to the original names. 

Currently, one example is included with the project, providing a basic guide to using a few functions provided within the domains component of the XML API client library. The example illustrates how to obtain a complete list of domains regardless of how many domains the account has (because the original the original library this was forked from could only list 20 domains maximum). This example is the same example provided below. 

As further progress is made to make this client library feature complete, additional examples will be included to illustrate that different major components of the client library.

#### Why does this fork exist?
**The library function names, and expected paramters and return variables are not frozen as the original library and are subject to change, be advised before including this library in your project, you may wish to create a fork or use a sp[ecifci commit if you decide touse it in a production application.** The reason for this is that this library was forked from a project that is not complete and provided no documentation informing developers which portions of the API were incorporated, any plan on how completion would occur, if the API was frozen and lastly the project maintainer was not responsive when issuing pull requests.

For example, the "DomainsGetList" function would only return your newest 20 domains with no way to make simple requests that included anything past the first page of results, and the pagination data was not even being cached in the APIResponse object. In addition, 

### Examples
The following example illustrates how one would obtain a list of Domains and request different chunks of the users registered domains using the pagination system built into the API but unfortunately not included in the original XML API Client.

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

