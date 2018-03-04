package namecheap

import (
	"errors"
	"fmt"
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
	TotalItems  int `xml:"TotalItems,attr"`
	CurrentPage int `xml:"CurrentPage,attr"`
	PageSize    int `xml:"PageSize,attr"`
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

func (client *Client) DomainsGetList(currentPage uint, pageSize uint) ([]DomainGetListResult, Paging, error) {
	if pageSize > 100 {
		// Maximum page size supported by the Namecheap API
		pageSize = 100
	}
	requestInfo := &ApiRequest{
		command: domainsGetList,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("CurrentPage", strconv.Itoa(int(currentPage)))
	requestInfo.params.Set("PageSize", strconv.Itoa(int(pageSize)))
	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp: ", resp)
	paging := Paging{
		TotalItems:  resp.params.Get("TotalItems"),
		CurrentPage: resp.params.Get("CurrentPage"),
		PageSize:    resp.params.Get("PageSize"),
	}
	fmt.Println("PAGING:")
	fmt.Println("Total Items  : ", paging.TotalItems)
	fmt.Println("Current Page : ", paging.CurrentPage)
	fmt.Println("Page Size    : ", paging.PageSize)

	return resp.Domains, paging, nil
}

func (client *Client) DomainGetInfo(domainName string) (*DomainInfo, error) {
	requestInfo := &ApiRequest{
		command: domainsGetInfo,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainName", domainName)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainInfo, nil
}

func (client *Client) DomainsCheck(domainNames ...string) ([]DomainCheckResult, error) {
	requestInfo := &ApiRequest{
		command: domainsCheck,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainList", strings.Join(domainNames, ","))
	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainsCheck, nil
}

func (client *Client) DomainsTLDList() ([]TLDListResult, error) {
	requestInfo := &ApiRequest{
		command: domainsTLDList,
		method:  "POST",
		params:  url.Values{},
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.TLDList, nil
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

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainCreate, nil
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
