package namecheap

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

const (
	domainsGetList = "namecheap.domains.getList"
	domainsGetInfo = "namecheap.domains.getInfo"
	domainsCheck   = "namecheap.domains.check"
	domainsCreate  = "namecheap.domains.create"
	domainsTLDList = "namecheap.domains.getTldList"
	domainsRenew   = "namecheap.domains.renew"
	// Domain `ListType` Filter
	// https://www.namecheap.com/support/api/methods/domains/get-list.aspx
	ALL      = "ALL"
	EXPIRING = "EXPIRING"
	EXPIRED  = "EXPIRED"
	// Domain `SortBy` Filter
	// https://www.namecheap.com/support/api/methods/domains/get-list.aspx
	NAME_ASC         = "NAME"
	NAME_DESC        = "NAME_DESC"
	EXPIRE_DATE_ASC  = "EXPIREDATE"
	EXPIRE_DATE_DESC = "EXPIREDATE_DESC"
	CREATE_DATE_ASC  = "CREATEDATE"
	CREATE_DATE_DESC = "CREATEDATE_DESC"
)

// DomainGetListResult represents the data returned by 'domains.getList'
type DomainGetListResult struct {
	ID         int    `xml:"ID,attr"`
	Name       string `xml:"Name,attr"`
	User       string `xml:"User,attr"`
	Created    string `xml:"Created,attr"`
	Expires    string `xml:"Expires,attr"`
	IsExpired  bool   `xml:"IsExpired,attr"`
	IsLocked   bool   `xml:"IsLocked,attr"`
	AutoRenew  bool   `xml:"AutoRenew,attr"`
	WhoisGuard string `xml:"WhoisGuard,attr"`
}

type Paging struct {
	TotalItems  uint `xml:"TotalItems,attr"`
	CurrentPage uint `xml:"CurrentPage,attr"`
	PageSize    uint `xml:"PageSize,attr"`
}

// DomainInfo represents the data returned by 'domains.getInfo'
type DomainInfo struct {
	ID         int        `xml:"ID,attr"`
	Name       string     `xml:"DomainName,attr"`
	Owner      string     `xml:"OwnerName,attr"`
	Created    string     `xml:"DomainDetails>CreatedDate"`
	Expires    string     `xml:"DomainDetails>ExpiredDate"`
	IsExpired  bool       `xml:"IsExpired,attr"`
	IsLocked   bool       `xml:"IsLocked,attr"`
	AutoRenew  bool       `xml:"AutoRenew,attr"`
	DNSDetails DNSDetails `xml:"DnsDetails"`
	Whoisguard Whoisguard `xml:"Whoisguard"`
}

type DNSDetails struct {
	ProviderType  string   `xml:"ProviderType,attr"`
	IsUsingOurDNS bool     `xml:"IsUsingOurDNS,attr"`
	Nameservers   []string `xml:"Nameserver"`
}

type Whoisguard struct {
	Enabled     bool   `xml:"Enabled,attr"`
	ID          int64  `xml:"ID"`
	ExpiredDate string `xml:"ExpiredDate"`
}

type DomainCheckResult struct {
	Domain                   string  `xml:"Domain,attr"`
	Available                bool    `xml:"Available,attr"`
	IsPremiumName            bool    `xml:"IsPremiumName,attr"`
	PremiumRegistrationPrice float64 `xml:"PremiumRegistrationPrice,attr"`
	PremiumRenewalPrice      float64 `xml:"PremiumRenewalPrice,attr"`
	PremiumRestorePrice      float64 `xml:"PremiumRestorePrice,attr"`
	PremiumTransferPrice     float64 `xml:"PremiumTransferPrice,attr"`
	IcannFee                 float64 `xml:"IcannFee,attr"`
}

type TLDListResult struct {
	Name string `xml:"Name,attr"`
}

type DomainCreateResult struct {
	Domain            string  `xml:"Domain,attr"`
	Registered        bool    `xml:"Registered,attr"`
	ChargedAmount     float64 `xml:"ChargedAmount,attr"`
	DomainID          int     `xml:"DomainID,attr"`
	OrderID           int     `xml:"OrderID,attr"`
	TransactionID     int     `xml:"TransactionID,attr"`
	WhoisguardEnable  bool    `xml:"WhoisguardEnable,attr"`
	NonRealTimeDomain bool    `xml:"NonRealTimeDomain,attr"`
}

type DomainRenewResult struct {
	DomainID      int     `xml:"DomainID,attr"`
	Name          string  `xml:"DomainName,attr"`
	Renewed       bool    `xml:"Renew,attr"`
	ChargedAmount float64 `xml:"ChargedAmount,attr"`
	OrderID       int     `xml:"OrderID,attr"`
	TransactionID int     `xml:"TransactionID,attr"`
	ExpireDate    string  `xml:"DomainDetails>ExpiredDate"`
}

type DomainCreateOption struct {
	AddFreeWhoisguard bool
	WGEnabled         bool
	Nameservers       []string
}

func (client *Client) DomainsGetList() (uint, error) {
	r, err := client.DomainsListAPIRequest(1, 1)
	if err != nil {
		return r.TotalItems, err
	}
	return r.TotalItems, err
}

// TODO: These function names are kinda awful, a overhaul of the library should address renaming these to give
// a more readable API and library usage that is intiutive, then alias to make it backwards compatible
func (client *Client) DomainsGetList(currentPage uint, pageSize uint) ([]DomainGetListResult, Paging, error) {
	r, err := client.DomainsListAPIRequest(currentPage, pageSize)
	p := Paging{
		TotalItems:  r.TotalItems,
		CurrentPage: r.CurrentPage,
		PageSize:    r.PageSize,
	}
	return r.Domains, p, err
}

func (client *Client) DomainsGetCompleteList() (domains []DomainGetListResult, err error) {
	r, err := client.DomainsListAPIRequest(1, maxPerPage)
	if err != nil {
		return nil, err
	}

	// Load first 100 items into map
	domains = append(domains, r.Domains...)
	if r.TotalItems > maxPerPage {
		remaining := (r.TotalItems - maxPerPage)
		quotient := (remaining / maxPerPage)
		if quotient != 0 {
			// Start from the second page because the initial apge is scrapped to get the
			// initial paging object and so +2 is added to quotient to request each page, and
			// an additonal +1 to request the remainder
			for currentPage := 2; uint(currentPage) < (quotient + 3); currentPage++ {
				r, err = client.DomainsListAPIRequest(uint(currentPage), maxPerPage)
				if err != nil {
					return domains, err
				}
				domains = append(domains, r.Domains...)
			}
		} else {
			r, err = client.DomainsListAPIRequest(2, maxPerPage)
			if err != nil {
				return domains, err
			}
			domains = append(domains, r.Domains...)
		}
	}
	return domains, nil
}

func (client *Client) DomainGetInfo(domainName string) (*DomainInfo, error) {
	requestInfo := &ApiRequest{
		command: domainsGetInfo,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)

	r, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}
	return r.DomainInfo, nil
}

func (client *Client) DomainsCheck(domainNames ...string) ([]DomainCheckResult, error) {
	requestInfo := &ApiRequest{
		command: domainsCheck,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainList", strings.Join(domainNames, ","))
	r, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return r.DomainsCheck, nil
}

func (client *Client) DomainsTLDList(currentPage int) ([]TLDListResult, Paging, error) {
	requestInfo := &ApiRequest{
		command: domainsTLDList,
		method:  "POST",
		params:  url.Values{},
	}

	r, err := client.do(requestInfo)
	if err != nil {
		return nil, Paging{}, err
	}
	return r.TLDList, Paging{TotalItems: r.TotalItems, CurrentPage: r.CurrentPage, PageSize: r.PageSize}, nil
}

func (client *Client) DomainCreate(domainName string, years int, options ...DomainCreateOption) (*DomainCreateResult, error) {
	if client.Registrant == nil {
		return nil, errors.New("Registrant information on client cannot be empty")
	}

	requestInfo := &ApiRequest{
		command: domainsCreate,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainName", domainName)
	requestInfo.params.Set("Years", strconv.Itoa(years))
	for _, opt := range options {
		if opt.AddFreeWhoisguard {
			requestInfo.params.Set("AddFreeWhoisguard", "yes")
		}
		if opt.WGEnabled {
			requestInfo.params.Set("WGEnabled", "yes")
		}
		if len(opt.Nameservers) > 0 {
			requestInfo.params.Set("Nameservers", strings.Join(opt.Nameservers, ","))
		}
	}
	if err := client.Registrant.addValues(requestInfo.params); err != nil {
		return nil, err
	}

	r, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return r.DomainCreate, nil
}

func (client *Client) DomainRenew(domainName string, years int) (*DomainRenewResult, error) {
	requestInfo := &ApiRequest{
		command: domainsRenew,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)
	requestInfo.params.Set("Years", strconv.Itoa(years))

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainRenew, nil
}
