// Package namecheap implements a client for the Namecheap API.
//
// In order to use this package you will need a Namecheap account and your API Token.
package namecheap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultBaseURL = "https://api.namecheap.com/xml.response"

	// VALIDATION
	validDomainCharacters = "abcdefghijklmnopqrstuvwxyz0123456789-"
	// Number of domains to be listed on a page. Minimum value is 10, and maximum value is 100.
	// Default Value: 20
	// https://www.namecheap.com/support/api/methods/domains/get-list.aspx
	minPerPage     = 10
	maxPerPage     = 100
	minCurrentPage = 1
	maxCurrentPage = 999
)

// Namecheap API Global Paramters
// https://www.namecheap.com/support/api/global-parameters.aspx

// A map may be the best datatype to store this information so it can be easily printed
// for the user
//=================================================================================================
// Name         Type    Max     Req     Description                                              //
//=================================================================================================
// ApiUser 	String 	20 	Yes 	Username required to access the API
// ApiKey 	String 	50 	Yes 	Password required used to access the API
// Command 	String 	80 	Yes 	Command for execution
// UserName 	String 	20 	Yes 	The Username on which a command is executed.Generally,
//                                      the values of ApiUser and UserName parameters are the same.
// ClientIp 	String 	15 	Yes 	IP address of the client accessing your application
//                                      (End-user IP address)
//

// Namecheap API Global Error Codes
// https://www.namecheap.com/support/api/global-parameters.aspx

// A map may be the best datatype to store this information so it can be easily printed
// for the user

//1010101 		Parameter APIUser is missing
//1030408 		Unsupported authentication type
//1010104 		Parameter Command is missing
//1010102, 1011102 	Parameter APIKey is missing
//1010105, 1011105 	Parameter ClientIP is missing
//1050900 		Unknown error when validating APIUser
//1011150 		Parameter RequestIP is invalid
//1017150 		Parameter RequestIP is disabled or locked
//1017105 		Parameter ClientIP is disabled or locked
//1017101 		Parameter ApiUser is disabled or locked
//1017410 		Too many declined payments
//1017411 		Too many login attempts
//1019103 		Parameter UserName is not available
//1016103 		Parameter UserName is unauthorized
//1017103 		Parameter UserName is disabled or locked

// Client represents a client used to make calls to the Namecheap API.
type Client struct {
	ApiUser    string
	ApiToken   string
	UserName   string
	HttpClient *http.Client

	// Base URL for API requests.
	// Defaults to the public Namecheap API,
	// but can be set to a different endpoint (e.g. the sandbox).
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	*Registrant
}

type ApiRequest struct {
	method  string
	command string
	params  url.Values
}

type ApiResponse struct {
	Status             string                    `xml:"Status,attr"`
	Command            string                    `xml:"RequestedCommand"`
	TLDList            []TLDListResult           `xml:"CommandResponse>Tlds>Tld"`
	Domains            []DomainGetListResult     `xml:"CommandResponse>DomainGetListResult>Domain"`
	DomainInfo         *DomainInfo               `xml:"CommandResponse>DomainGetInfoResult"`
	DomainDNSHosts     *DomainDNSGetHostsResult  `xml:"CommandResponse>DomainDNSGetHostsResult"`
	DomainDNSSetHosts  *DomainDNSSetHostsResult  `xml:"CommandResponse>DomainDNSSetHostsResult"`
	DomainCreate       *DomainCreateResult       `xml:"CommandResponse>DomainCreateResult"`
	DomainRenew        *DomainRenewResult        `xml:"CommandResponse>DomainRenewResult"`
	DomainsCheck       []DomainCheckResult       `xml:"CommandResponse>DomainCheckResult"`
	DomainNSInfo       *DomainNSInfoResult       `xml:"CommandResponse>DomainNSInfoResult"`
	DomainDNSSetCustom *DomainDNSSetCustomResult `xml:"CommandResponse>DomainDNSSetCustomResult"`
	UsersGetPricing    []UsersGetPricingResult   `xml:"CommandResponse>UserGetPricingResult>ProductType"`
	WhoisguardList     []WhoisguardGetListResult `xml:"CommandResponse>WhoisguardGetListResult>Whoisguard"`
	WhoisguardEnable   whoisguardEnableResult    `xml:"CommandResponse>WhoisguardEnableResult"`
	WhoisguardDisable  whoisguardDisableResult   `xml:"CommandResponse>WhoisguardDisableResult"`
	WhoisguardRenew    *WhoisguardRenewResult    `xml:"CommandResponse>WhoisguardRenewResult"`
	TotalItems         uint                      `xml:"CommandResponse>Paging>TotalItems"`
	CurrentPage        uint                      `xml:"CommandResponse>Paging>CurrentPage"`
	PageSize           uint                      `xml:"CommandResponse>Paging>PageSize"`

	Errors ApiErrors `xml:"Errors>Error"`
}

// ApiError is the format of the error returned in the api responses.
type ApiError struct {
	Number  int    `xml:"Number,attr"`
	Message string `xml:",innerxml"`
}

func (err *ApiError) Error() string {
	return err.Message
}

// ApiErrors holds multiple ApiError's but implements the error interface
type ApiErrors []ApiError

func (errs ApiErrors) Error() string {
	errMsg := ""
	for _, apiError := range errs {
		errMsg += fmt.Sprintf("Error %d: %s\n", apiError.Number, apiError.Message)
	}
	return errMsg
}

// VALIDATION
func ValidDomainName(d string) bool {
	for _, char := range d {
		if !strings.Contains(validDomainCharacters, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

func ValidatePageSize(pageSize uint) {
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	} else if pageSize < minPerPage {
		pageSize = minPerPage
	}
}

func ValidateCurrentPage(page uint) {
	if page > maxCurrentPage {
		page = maxCurrentPage
	} else if page < minCurrentPage {
		page = minCurrentPage
	}
}

func ValidateSearchTerm(searchTerm string) error {
	if len(searchTerm) <= 1 {
		searchTerm = ""
	} else if len(searchTerm) >= 128 {
		searchTerm = searchTerm[:128]
	}
	if ValidDomainName(searchTerm) {
		return errors.New("invalid domain characters in search term")
	}
	return nil
}

func ValidateListType(listType string) {
	if listType != ALL || listType != EXPIRING || listType != EXPIRED {
		listType = ALL
	}
}

func ValidateSortBy(sortBy string) {
	if sortBy != NAME_ASC || sortBy != NAME_DESC || sortBy != EXPIRE_DATE_ASC || sortBy != EXPIRE_DATE_DESC || sortBy != CREATE_DATE_ASC || sortBy != CREATE_DATE_DESC {
		sortBy = NAME_ASC
	}
}

// API CLIENT
func NewClient(apiUser, apiToken, userName string) *Client {
	return &Client{
		ApiUser:    apiUser,
		ApiToken:   apiToken,
		UserName:   userName,
		HttpClient: http.DefaultClient,
		BaseURL:    defaultBaseURL,
	}
}

// NewRegistrant associates a new registrant with the
func (client *Client) NewRegistrant(
	firstName, lastName,
	addr1, addr2,
	city, state, postalCode, country,
	phone, email string,
) {
	client.Registrant = newRegistrant(
		firstName, lastName,
		addr1, addr2,
		city, state, postalCode, country,
		phone, email,
	)
}

func (client *Client) do(request *ApiRequest) (*ApiResponse, error) {
	if request.method == "" {
		return nil, errors.New("request method cannot be blank")
	}

	body, _, err := client.sendRequest(request)
	if err != nil {
		return nil, err
	}

	resp := new(ApiResponse)
	if err = xml.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	if resp.Status == "ERROR" {
		return nil, resp.Errors
	}

	return resp, nil
}

func (client *Client) makeRequest(request *ApiRequest) (*http.Request, error) {
	p := request.params
	p.Set("ApiUser", client.ApiUser)
	p.Set("ApiKey", client.ApiToken)
	p.Set("UserName", client.UserName)
	// This param is required by the API, but not actually used.
	p.Set("ClientIp", "127.0.0.1")
	p.Set("Command", request.command)

	b := p.Encode()
	req, err := http.NewRequest(request.method, client.BaseURL, strings.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return req, nil
}

func (client *Client) sendRequest(request *ApiRequest) ([]byte, int, error) {
	req, err := client.makeRequest(request)
	if err != nil {
		return nil, 0, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return buf, resp.StatusCode, nil
}

func (client *Client) DomainsListAPIRequest(page uint, pageSize uint, searchTerm, listType, sortBy string) (*ApiResponse, error) {
	// VALIDATION
	// [pageSize] must be equal or GREATER than 10
	// [pageSize] must be equal or LESS than 100
	ValidatePageSize(pageSize)
	// [page] must be equal or GREATER than 1
	// [page] must be qual or LESS than 999 (sanity)
	ValidateCurrentPage(page)
	// [searchTerm] must be alphanumeric
	// [searchTerm] must have a length GREATER than 1
	// [searchTerm] must have a length LESS than 128
	ValidateSearchTerm(searchTerm)
	// [listType] can only be ALL, EXPIRING, or EXPIRED (Default: ALL)
	ValidateListType(listType)
	// [sortBy] can only be NAME, NAME_DESC, EXPIREDATE, EXPIREDATE_DESC, CREATEDATE, CREATEDATE_DESC
	ValidateSortBy(sortBy)

	requestInfo := &ApiRequest{
		command: domainsGetList,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("page", strconv.Itoa(int(page)))
	requestInfo.params.Set("pageSize", strconv.Itoa(int(pageSize)))

	r, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}
	if &r == nil {
		return nil, errors.New("Request struct failed to be assigned")
	} else if &r.TotalItems == nil {
		return nil, errors.New("Request fields related to Paging, specifically 'TotalItems', 'CurrentPage', and 'PageSize' failed to be assigned")
	}
	return r, nil
}
