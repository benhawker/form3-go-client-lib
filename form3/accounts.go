package form3

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	accountsPath string = "/organisation/accounts"
)

// Account represents a bank account that is registered with Form3. It is used to validate and allocate inbound payments.
type Account struct {
	Attributes     AccountAttributes `json:"attributes"`
	ID             string            `json:"id"`
	OrganisationID string            `json:"organisation_id"`
	Type           string            `json:"type"`
	Version        int               `json:"version"`
}

// AccountAttributes represents attributes of an Account
type AccountAttributes struct {
	Country                     string   `json:"country"`
	BaseCurrency                string   `json:"base_currency,omitempty"`
	AccountNumber               string   `json:"account_number,omitempty"`
	BankID                      string   `json:"bank_id,omitempty"`
	BankIDCode                  string   `json:"bank_id_code,omitempty"`
	Bic                         string   `json:"bic,omitempty"`
	Iban                        string   `json:"iban,omitempty"`
	Title                       string   `json:"title,omitempty"`
	FirstName                   string   `json:"first_name,omitempty"`
	BankAccountName             string   `json:"bank_account_name,omitempty"`
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names,omitempty"`
	AccountClassification       string   `json:"account_classification,omitempty"`
	JointAccount                bool     `json:"joint_account,omitempty"`
	AccountMatchingOptOut       bool     `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification     string   `json:"secondary_identification,omitempty"`
}

// Links -> represents the related links to returned resource(s)
//
// See https://api-docs.form3.tech/api.html?http#introduction-and-api-conventions-message-body-structure
// This follows the HATEOAS convention (https://en.wikipedia.org/wiki/HATEOAS)
type Links struct {
	First *string `json:"first,omitempty"`
	Last  *string `json:"last,omitempty"`
	Next  *string `json:"next,omitempty"`
	Prev  *string `json:"prev,omitempty"`
	Self  *string `json:"self"`
}

// AccountsService implements a service to manage accounts
// See https://api-docs.form3.tech/api.html?http#organisation-accounts
type AccountsService struct {
	client     *Client
	pagination Pagination
}

// NewAccountsService creates a new AccountsService.
func NewAccountsService(client *Client) *AccountsService {
	builder := &AccountsService{
		client:     client,
		pagination: NewPagination(),
	}
	return builder
}

type fetchAccountAPIResponse struct {
	Data  Account `json:"data"`
	Links Links   `json:"links"`
}

type listAccountsAPIResponse struct {
	Data  []Account `json:"data"`
	Links Links     `json:"links"`
}

type createAccountsAPIPayload struct {
	Data Account `json:"data"`
}

// Fetch -> Get a single account using the account ID.
//
// GET /v1/organisation/accounts/{account_id}
func (s *AccountsService) Fetch(ctx context.Context, id string) (*Account, *http.Response, error) {
	res, err := s.client.MakeRequest(ctx, MakeRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", accountsPath, id),
	})
	if err != nil {
		return nil, res, err
	}

	var ret fetchAccountAPIResponse
	if err := s.client.Decode(res, &ret); err != nil {
		return nil, res, err
	}

	return &ret.Data, res, nil
}

// List -> List accounts with the ability to filter and page.
// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
//
// GET /v1/organisation/accounts?page[number]={page_number}&page[size]={page_size}&filter[{attribute}]={filter_value}
func (s *AccountsService) List(ctx context.Context) ([]Account, *http.Response, error) {
	res, err := s.client.MakeRequest(ctx, MakeRequestOptions{
		Method: "GET",
		Path:   accountsPath,
		Params: s.pagination.Params(),
	})
	if err != nil {
		return nil, nil, err
	}

	var ret listAccountsAPIResponse
	if err := s.client.Decode(res, &ret); err != nil {
		return nil, nil, err
	}
	return ret.Data, res, nil
}

// Create -> Register an existing bank account with Form3 or create a new one.
//
// POST /v1/organisation/accounts
//
// The country attribute must be specified as a minimum.
// Depending on the country, other attributes such as bank_id and bic are mandatory.
//
// Form3 generates account numbers and IBANs, where appropriate, in the following cases:
// - If no account number or IBAN is provided, Form3 generates a valid account number (see below). If supported by the country, an IBAN is also generated.
// - If an account number is provided but the IBAN is empty, Form3 generates an IBAN if supported by the country.
// - If only an IBAN is provided, the account number will be left empty.
// - Note that a given bank_id and bic need to be registered with Form3 and connected to your organisation ID.
// See https://api-docs.form3.tech/api.html?shell#organisation-accounts-create for further details.
func (s *AccountsService) Create(ctx context.Context, account *Account) (*Account, *http.Response, error) {
	data := &createAccountsAPIPayload{Data: *account}

	res, err := s.client.MakeRequest(ctx, MakeRequestOptions{
		Method: "POST",
		Path:   accountsPath,
		Body:   data,
	})
	if err != nil {
		return nil, res, err
	}

	var ret fetchAccountAPIResponse
	if err := s.client.Decode(res, &ret); err != nil {
		return nil, res, err
	}

	return &ret.Data, res, nil
}

// Delete -> Delete an account
//
// DELETE /v1/organisation/accounts/:id?version=:version
//
// No response body returned.
// Potential status codes:
// - 204	No Content	Resource has been successfully deleted
// - 404	Not Found	Specified resource does not exist
// - 409	Conflict	Specified version incorrect
func (s *AccountsService) Delete(ctx context.Context, id string, version int) (bool, *http.Response, error) {
	params := url.Values{}
	params.Add("version", strconv.Itoa(version))

	res, err := s.client.MakeRequest(ctx, MakeRequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s", accountsPath, id),
		Params: params,
	})
	if err != nil {
		return false, res, err
	}

	return true, res, nil
}

// Number -> page number requested. Defaults to 0.
func (s *AccountsService) Number(number int) *AccountsService {
	s.pagination.Number = number
	return s
}

// Size -> size is the max number of resources to return. Defaults to 10.
func (s *AccountsService) Size(size int) *AccountsService {
	s.pagination.Size = size
	return s
}
