package form3

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_FetchAccount_Success(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts/158f775c-4ecd-4861-b33d-30df9a29de78", http.StatusOK, accountJSON)
	defer srv.Close()

	account, res, err := client.Accounts().Fetch(context.Background(), "158f775c-4ecd-4861-b33d-30df9a29de78")
	if err != nil {
		t.Error(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if account.ID != "158f775c-4ecd-4861-b33d-30df9a29de78" {
		t.Error("Expected: 158f775c-4ecd-4861-b33d-30df9a29de78", "Got:", account.ID)
	}
}
func Test_FetchAccount_NotFound_Failure(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts/not-an-id", 404, "")
	defer srv.Close()

	account, res, _ := client.Accounts().Fetch(context.TODO(), "not-an-id")

	if account != nil {
		t.Errorf("Expected account to be nil but got %v", account)
	}

	if http.StatusNotFound != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}
}

func Test_ListAccounts_Success(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts", http.StatusOK, accountsJSON)
	defer srv.Close()

	accounts, res, err := client.Accounts().List(context.Background())
	if err != nil {
		t.Error(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if len(accounts) != 2 {
		t.Error("Expected:", 2, "Got:", len(accounts))
	}

	fetchedAccountIDs := make([]string, len(accounts))
	for index, account := range accounts {
		fetchedAccountIDs[index] = account.ID
	}

	expectedIDs := []string{"158f775c-4ecd-4861-b33d-30df9a29de78", "258f775c-4ecd-4861-b33d-30df9a29de78"}
	if !reflect.DeepEqual(fetchedAccountIDs, expectedIDs) {
		t.Error("Expected:", expectedIDs, "Got:", fetchedAccountIDs)
	}
}

func Test_ListAccountsWithPagination_Success(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts", http.StatusOK, accountsJSON)
	defer srv.Close()

	accounts, res, err := client.Accounts().Number(1).Size(10).List(context.Background())
	if err != nil {
		t.Error(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if len(accounts) != 2 {
		t.Error("Expected:", 2, "Got:", len(accounts))
	}

	fetchedAccountIDs := make([]string, len(accounts))
	for index, account := range accounts {
		fetchedAccountIDs[index] = account.ID
	}

	expectedIDs := []string{"158f775c-4ecd-4861-b33d-30df9a29de78", "258f775c-4ecd-4861-b33d-30df9a29de78"}
	if !reflect.DeepEqual(fetchedAccountIDs, expectedIDs) {
		t.Error("Expected:", expectedIDs, "Got:", fetchedAccountIDs)
	}
}

func Test_CreateAccount_Success(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts", http.StatusOK, accountJSON)
	defer srv.Close()

	account, res, err := client.Accounts().Create(context.Background(), &Account{})
	if err != nil {
		t.Error(err)
	}

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if account.ID != "158f775c-4ecd-4861-b33d-30df9a29de78" {
		t.Error("Expected: 158f775c-4ecd-4861-b33d-30df9a29de78", "Got:", account.ID)
	}
}

func Test_CreateAccount_InvalidPayload_Failure(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts", http.StatusBadRequest, invalidPayload)
	defer srv.Close()

	_, res, err := client.Accounts().Create(context.Background(), &Account{})
	if err == nil {
		t.Error("Expected: error", "Got: nil")
	}

	if http.StatusBadRequest != res.StatusCode {
		t.Error("Expected:", http.StatusBadRequest, "Got:", res.StatusCode)
	}
}

func Test_DeleteAccount_Success(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts/158f775c-4ecd-4861-b33d-30df9a29de78", http.StatusNoContent, "")
	defer srv.Close()

	ok, res, err := client.Accounts().Delete(context.Background(), "158f775c-4ecd-4861-b33d-30df9a29de78", 0)
	if err != nil {
		t.Error(err)
	}

	if ok != true {
		t.Error("Expected: true Got:", ok)
	}

	if http.StatusNoContent != res.StatusCode {
		t.Error("Expected:", http.StatusNoContent, "Got:", res.StatusCode)
	}
}

func Test_DeleteAccount_NotFound_Failure(t *testing.T) {
	client, srv := testClient("/v1/organisation/accounts/not-found-uuid", http.StatusNotFound, "")
	defer srv.Close()

	ok, res, err := client.Accounts().Delete(context.Background(), "not-found-uuid", 0)
	if err == nil {
		t.Error("Expected: error", "Got: nil")
	}

	if ok != false {
		t.Error("Expected: false Got:", ok)
	}

	if http.StatusNotFound != res.StatusCode {
		t.Error("Expected:", http.StatusNotFound, "Got:", res.StatusCode)
	}
}

func serverMock(path string, handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(path, handler)

	srv := httptest.NewServer(mux)

	return srv
}

func testClient(path string, statusCode int, responseBody string) (*Client, *httptest.Server) {
	srv := serverMock(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	})

	// httptest.Server returns url as string (https://golang.org/pkg/net/http/httptest/#Server)
	u, _ := url.Parse(srv.URL)

	client, _ := NewClient(
		SetScheme(u.Scheme),
		SetHost(u.Host),
	)

	return client, srv
}

var accountJSON = `{
    "data": {
        "type": "accounts",
        "id": "158f775c-4ecd-4861-b33d-30df9a29de78",
        "version": 0,
        "organisation_id": "158f775c-4ecd-4861-b33d-30df9a29de78",
        "attributes": {
            "country": "GB",
            "base_currency": "GBP",
            "account_number": "1112223",
            "bank_id": "1234567",
            "bank_id_code": "GBABC"
        }
    }
}`

var accountsJSON = `{
	"data": [{
		"attributes": {
			"account_number": "11122233",
			"alternative_bank_account_names": null,
			"bank_id": "111222",
			"bank_id_code": "GBABC",
			"base_currency": "GBP",
			"bic": "BUKBGB22",
			"country": "GB",
			"iban": "GB1400080001001234567890"
		},
		"created_on": "2020-06-30T15:16:30.270Z",
		"id": "158f775c-4ecd-4861-b33d-30df9a29de78",
		"modified_on": "2020-06-30T15:16:30.270Z",
		"organisation_id": "158f775d-4ecd-4861-b33d-30df9a29de78",
		"type": "accounts",
		"version": 0
	},
	{
		"attributes": {
			"account_number": "11122239",
			"alternative_bank_account_names": null,
			"bank_id": "111229",
			"bank_id_code": "GBABC",
			"base_currency": "GBP",
			"bic": "BUKBGB22",
			"country": "GB",
			"iban": "GB1400080001001234567899"
		},
		"created_on": "2020-06-30T15:16:30.270Z",
		"id": "258f775c-4ecd-4861-b33d-30df9a29de78",
		"modified_on": "2020-06-30T15:16:30.270Z",
		"organisation_id": "158f775d-4ecd-4861-b33d-30df9a29de78",
		"type": "accounts",
		"version": 0
	}],
	"links": {
		"first": "/v1/organisation/accounts?page%5Bnumber%5D=first&page%5Bsize%5D=1",
		"last": "/v1/organisation/accounts?page%5Bnumber%5D=last&page%5Bsize%5D=1",
		"next": "/v1/organisation/accounts?page%5Bnumber%5D=1&page%5Bsize%5D=1",
		"self": "/v1/organisation/accounts?page%5Bnumber%5D=0&page%5Bsize%5D=1"
	}
}`

var invalidPayload = `{
	"data": {
		}
	}
}`
