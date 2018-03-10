# Namecheap Go XML API Client
A Go client library that aims to be consistent, provide complete functionality (it is not currently complete), and additionally provide a thin abstraction layer for improving ease-of-use and access to the Namecheap XML [API](https://www.namecheap.com/support/api/intro.aspx).

A minimalistic Namecheap API that is not currently yet feature complete and there is not yet documentation detailing feature completeness or a roadmap outlining a path to feature completeness. The library in additional to providing a complete implementation of the XML API, it will provide a thin abstraction layer to simplify use, abstract additional features to make using the API more intuitive. As the library evolves, the function names will be altered to make use more intuitive and code using the library more readable but the original functions names will still be accessible and usable by aliasing the altered names to the original names. 

Currently, one example is included with the project, providing a basic guide to using a few functions provided within the domains component of the XML API client library. The example illustrates how to obtain a complete list of domains regardless of how many domains the account has (because the original the original library this was forked from could only list 20 domains maximum). This example is the same example provided below. 

As further progress is made to make this client library feature complete, additional examples will be included to illustrate that different major components of the client library.

#### Why does this fork exist?
**The library function names, and expected paramters and return variables are not frozen as the original library and are subject to change, be advised before including this library in your project, you may wish to create a fork or use a sp[ecifci commit if you decide touse it in a production application.** The reason for this is that this library was forked from a project that is not complete and provided no documentation informing developers which portions of the API were incorporated, any plan on how completion would occur, if the API was frozen and lastly the project maintainer was not responsive when issuing pull requests.

For example, the "DomainsGetList" function would only return your newest 20 domains with no way to make simple requests that included anything past the first page of results, and the pagination data was not even being cached in the APIResponse object. In addition, 

### Namecheap API
Below are a list of the API endpoints provided by Namecheap, over development process this will be used to both indicate which functionality has bee completed , what is not and serve as the foundation for the roadmap to API client completion.

**domains**

    getList — Returns a list of domains for the particular user.
    getContacts — Gets contact information of the requested domain.
    create — Registers a new domain name.
    getTldList — Returns a list of tlds
    setContacts — Sets contact information for the domain.
    check — Checks the availability of domains.
    reactivate — Reactivates an expired domain.
    renew — Renews an expiring domain.
    getRegistrarLock — Gets the RegistrarLock status of the requested domain.
    setRegistrarLock — Sets the RegistrarLock status for a domain.
    getInfo — Returns information about the requested domain. 

**domains.dns**

    setDefault — Sets domain to use our default DNS servers. Required for free services like Host record management, URL forwarding, email forwarding, dynamic dns and other value added services.
    setCustom — Sets domain to use custom DNS servers. NOTE: Services like URL forwarding, Email forwarding, Dynamic DNS will not work for domains using custom nameservers.
    getList — Gets a list of DNS servers associated with the requested domain.
    getHosts — Retrieves DNS host record settings for the requested domain.
    getEmailForwarding — Gets email forwarding settings for the requested domain
    setEmailForwarding — Sets email forwarding for a domain name.
    setHosts — Sets DNS host records settings for the requested domain. 

**domains.ns**

    create — Creates a new nameserver.
    delete — Deletes a nameserver associated with the requested domain.
    getInfo — Retrieves information about a registered nameserver.
    update — Updates the IP address of a registered nameserver. 

**domains.transfer**

    create — Transfers a domain to Namecheap. You can only transfer .biz, .ca, .cc, .co, .co.uk, .com, .com.es, .com.pe, .es, .in, .info, .me, .me.uk, .mobi, .net, .net.pe, .nom.es, .org, .org.es, .org.pe, .org.uk, .pe, .tv, .us domains through API at this time.
    getStatus — Gets the status of a particular transfer.
    updateStatus — Updates the status of a particular transfer. Allows you to re-submit the transfer after releasing the registry lock.
    getList — Gets the list of domain transfers. 

**ssl**

    create — Creates a new SSL certificate.
    getList — Returns a list of SSL certificates for the particular user.
    parseCSR — Parsers the CSR
    getApproverEmailList — Gets approver email list for the requested certificate.
    activate — Activates a newly purchased SSL certificate.
    resendApproverEmail — Resends the approver email.
    getInfo — Retrieves information about the requested SSL certificate
    renew — Renews an SSL certificate.
    reissue — Reissues an SSL certificate.
    resendfulfillmentemail — Resends the fulfilment email containing the certificate.
    purchasemoresans — Purchases more add-on domains for already purchased certificate.
    revokecertificate — Revokes a re-issued SSL certificate.
    editDCVMethod — Sets new domain control validation (DCV) method for a certificate or serves as 'retry' mechanism

**users**

    getPricing — Returns pricing information for a requested product type.
    getBalances — Gets information about fund in the user's account.This method returns the following information: Available Balance, Account Balance, Earned Amount, Withdrawable Amount and Funds Required for AutoRenew.
    changePassword — Changes password of the particular user's account.
    update — Updates user account information for the particular user.
    createaddfundsrequest — Creates a request to add funds through a credit card
    getAddFundsStatus — Gets the status of add funds request.
    create — Creates a new account at NameCheap under this ApiUser.
    login — Validates the username and password of user accounts you have created using the API command namecheap.users.create.
    resetPassword — When you call this API, a link to reset password will be emailed to the end user's profile email id.The end user needs to click on the link to reset password. 

**users.address**

    create — Creates a new address for the user
    delete — Deletes the particular address for the user.
    getInfo — Gets information for the requested addressID.
    getList — Gets a list of addressIDs and addressnames associated with the user account.
    setDefault — Sets default address for the user.
    update — Updates the particular address of the user 

**whoisguard**

    changeemailaddress — Changes WhoisGuard email address
    enable — Enables WhoisGuard privacy protection.
    disable — Disables WhoisGuard privacy protection.
    unallot — Unallots WhoisGuard privacy protection.
    discard — Discards whoisguard.
    allot — Allots WhoisGuard
    getList — Gets the list of WhoisGuard privacy protection.
    renew — Renews WhoisGuard privacy protection.

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

