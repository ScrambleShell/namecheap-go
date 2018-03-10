# Namecheap Go XML API Client
A Go client library that aims to be consistent, provide complete functionality (it is not currently complete), and additionally provide a thin abstraction layer for improving ease-of-use and access to the Namecheap XML [API](https://www.namecheap.com/support/api/intro.aspx).

A minimalistic Namecheap API that is not currently yet feature complete and there is not yet documentation detailing feature completeness or a roadmap outlining a path to feature completeness. The library in additional to providing a complete implementation of the XML API, it will provide a thin abstraction layer to simplify use, abstract additional features to make using the API more intuitive. As the library evolves, the function names will be altered to make use more intuitive and code using the library more readable but the original functions names will still be accessible and usable by aliasing the altered names to the original names. 

Currently, one example is included with the project, providing a basic guide to using a few functions provided within the domains component of the XML API client library. The example illustrates how to obtain a complete list of domains regardless of how many domains the account has (because the original the original library this was forked from could only list 20 domains maximum). This example is the same example provided below. 

As further progress is made to make this client library feature complete, additional examples will be included to illustrate that different major components of the client library.

#### Why does this fork exist?
**The library function names, and expected paramters and return variables are not frozen as the original library and are subject to change, be advised before including this library in your project, you may wish to create a fork or use a sp[ecifci commit if you decide touse it in a production application.** The reason for this is that this library was forked from a project that is not complete and provided no documentation informing developers which portions of the API were incorporated, any plan on how completion would occur, if the API was frozen and lastly the project maintainer was not responsive when issuing pull requests.

For example, the "DomainsGetList" function would only return your newest 20 domains with no way to make simple requests that included anything past the first page of results, and the pagination data was not even being cached in the APIResponse object. In addition, 
### Development Roadmap
#### Namecheap API
Below are a list of the API endpoints provided by Namecheap, over development process this will be used to both indicate which functionality has bee completed , what is not and serve as the foundation for the roadmap to API client completion.

##### Domains
The first priority for development is access to domains, listing domains, retrieving their details, registering domains, and the other domain related actions.

At the time of forking this component of the API was not complete, since forking, `Paging` has been added to allow developers to request every domain (previously the limit was 20). But the `SortBy`, `ListType`, and `SearchTerm` filters have not yet been added. See the [Namecheap API Documentation](https://www.namecheap.com/support/api/methods/domains/get-list.aspx) for details. 

  _getList_          — Returns a list of domains for the particular user.

The `getList` functions in the current Namecheap API 

  _getContacts_      — Gets contact information of the requested domain.
  _create_           — Registers a new domain name.
  _getTldList_       — Returns a list of tlds
  _setContacts_      — Sets contact information for the domain.
  _check_            — Checks the availability of domains.
  _reactivate_       — Reactivates an expired domain.
  _renew_            — Renews an expiring domain.
  _getRegistrarLock_ — Gets the RegistrarLock status of the requested domain.
  _setRegistrarLock_ — Sets the RegistrarLock status for a domain.
  _getInfo_          — Returns information about the requested domain. 

**domains.dns**
    _setDefault_         — Sets domain to use our default DNS servers. Required for free services like Host record management, URL forwarding, email forwarding, dynamic dns and other value added services.
    _setCustom_          — Sets domain to use custom DNS servers. NOTE: Services like URL forwarding, Email forwarding, Dynamic DNS will not work for domains using custom nameservers.
    _getList_            — Gets a list of DNS servers associated with the requested domain.
    _getHosts_           — Retrieves DNS host record settings for the requested domain.
    _getEmailForwarding_ — Gets email forwarding settings for the requested domain
    _setEmailForwarding_ — Sets email forwarding for a domain name.
    _setHosts_           — Sets DNS host records settings for the requested domain. 

**domains.ns**
    _create_  — Creates a new nameserver.
    _delete_  — Deletes a nameserver associated with the requested domain.
    _getInfo_ — Retrieves information about a registered nameserver.
    _update_  — Updates the IP address of a registered nameserver. 

**domains.transfer**
    _create_       — Transfers a domain to Namecheap. You can only transfer .biz, .ca, .cc, .co, .co.uk, .com, .com.es, .com.pe, .es, .in, .info, .me, .me.uk, .mobi, .net, .net.pe, .nom.es, .org, .org.es, .org.pe, .org.uk, .pe, .tv, .us domains through API at this time.
    _getStatus_    — Gets the status of a particular transfer.
    _updateStatus_ — Updates the status of a particular transfer. Allows you to re-submit the transfer after releasing the registry lock.
    _getList_      — Gets the list of domain transfers. 

~~**ssl**~~
*(Unncessary to implement, develoeprs should be encouraged to use ACME/Let's encrypt instead. It is a free SSL certificate provider, with easy automated certificate generation provided by the Mozilla project.)*
    ~~create~~                 — Creates a new SSL certificate.
    ~~getList~~                — Returns a list of SSL certificates for the particular user.
    ~~parseCSR~~               — Parsers the CSR
    ~~getApproverEmailList~~   — Gets approver email list for the requested certificate.
    ~~activate~~               — Activates a newly purchased SSL certificate.
    ~~resendApproverEmail~~    — Resends the approver email.
    ~~getInfo~~                — Retrieves information about the requested SSL certificate
    ~~renew~~                  — Renews an SSL certificate.
    ~~reissue~~                — Reissues an SSL certificate.
    ~~resendfulfillmentemail~~ — Resends the fulfilment email containing the certificate.
    ~~purchasemoresans~~       — Purchases more add-on domains for already purchased certificate.
    ~~revokecertificate~~      — Revokes a re-issued SSL certificate.
    ~~editDCVMethod~~          — Sets new domain control validation (DCV) method for a certificate or serves as 'retry' mechanism

**users**
    _getPricing_            — Returns pricing information for a requested product type.
    _getBalances_           — Gets information about fund in the user's account.This method returns the following information: Available Balance, Account Balance, Earned Amount, Withdrawable Amount and Funds Required for AutoRenew.
    _changePassword_        — Changes password of the particular user's account.
    _update_                — Updates user account information for the particular user.
    _createaddfundsrequest_ — Creates a request to add funds through a credit card
    _getAddFundsStatus_     — Gets the status of add funds request.
    _create_                — Creates a new account at NameCheap under this ApiUser.
    _login_                 — Validates the username and password of user accounts you have created using the API command namecheap.users.create.
    _resetPassword_         — When you call this API, a link to reset password will be emailed to the end user's profile email id.The end user needs to click on the link to reset password. 

**users.address**
    _create_     — Creates a new address for the user
    _delete_     — Deletes the particular address for the user.
    _getInfo_    — Gets information for the requested addressID.
    _getList_    — Gets a list of addressIDs and addressnames associated with the user account.
    _setDefault_ — Sets default address for the user.
    _update_     — Updates the particular address of the user 

**whoisguard**
    _changeemailaddress_ — Changes WhoisGuard email address
    _enable_             — Enables WhoisGuard privacy protection.
    _disable_            — Disables WhoisGuard privacy protection.
    _unallot_            — Unallots WhoisGuard privacy protection.
    _discard_            — Discards whoisguard.
    _allot_              — Allots WhoisGuard
    _getList_            — Gets the list of WhoisGuard privacy protection.
    _renew_              — Renews WhoisGuard privacy protection.

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

