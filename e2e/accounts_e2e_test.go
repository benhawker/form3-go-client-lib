package e2e

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	form3 "github.com/form3tech-oss/interview-accountapi/form3"
)

func Test_Run_e2e(t *testing.T) {

	// Init a client
	client, err := form3.NewClient(
		form3.SetScheme("http"),
	)
	if err != nil {
		t.Error("Error initializing client")
	}

	// Ensure we clean slate records prior to running
	// This approach is brittle and should be improved given more time.
	accounts, _, _ := client.Accounts().List(context.Background())

	if len(accounts) != 0 {
		t.Error("Records already present! Please clean slate the env prior to running these integration tests.")
		return
	}

	// An improved testing approach would generate these dynamically.
	uuid := "158f775c-4ecd-4861-b33d-30df9a29de78"
	uuid2 := "258f775c-4ecd-4861-b33d-30df9a29de78"

	// Create a valid account
	newAccount := generateAccount(uuid)
	account, res, err := client.Accounts().Create(context.Background(), newAccount)
	if err != nil {
		t.Error(err)
		t.Error("Error creating account")
	}

	if reflect.TypeOf(account).Elem().Name() != "Account" {
		t.Error("Expected:", "Account", "Got:", reflect.TypeOf(account).Elem().Name())
	}

	if http.StatusCreated != res.StatusCode {
		t.Error("Expected:", http.StatusCreated, "Got:", res.StatusCode)
	}

	if account.ID != newAccount.ID {
		t.Error("Expected:", newAccount.ID, "Got:", account.ID)
	}

	// Fetch it
	account, res, err = client.Accounts().Fetch(context.Background(), uuid)

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if account.ID != newAccount.ID {
		t.Error("Expected:", newAccount.ID, "Got:", account.ID)
	}

	// Create another valid account
	newAccount2 := generateAccount(uuid2)
	account, res, err = client.Accounts().Create(context.Background(), newAccount2)
	if err != nil {
		t.Error("Error creating account", err)
	}

	if reflect.TypeOf(account).Elem().Name() != "Account" {
		t.Error("Expected:", "Account", "Got:", reflect.TypeOf(account).Elem().Name())
	}

	if http.StatusCreated != res.StatusCode {
		t.Error("Expected:", http.StatusCreated, "Got:", res.StatusCode)
	}

	if account.ID != newAccount2.ID {
		t.Error("Expected:", newAccount2.ID, "Got:", account.ID)
	}

	// List accounts (expect 2)
	accounts, res, err = client.Accounts().List(context.Background())

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if len(accounts) != 2 {
		t.Error("Expected: 2 accounts", "Got:", len(accounts))
	}

	// Atrempt to create an account that will return a 409 Conflict
	account, res, err = client.Accounts().Create(context.Background(), newAccount2)
	if err == nil {
		t.Error("Expected error whilst creating account but it was nil")
	}

	if http.StatusConflict != res.StatusCode {
		t.Error("Expected:", http.StatusConflict, "Got:", res.StatusCode)
	}

	// List accounts (expect 2)
	accounts, res, err = client.Accounts().List(context.Background())

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if len(accounts) != 2 {
		t.Error("Expected: 2 accounts", "Got:", len(accounts))
	}

	// Delete 1 account
	ok, _, err := client.Accounts().Delete(context.Background(), uuid2, 0)

	if ok != true {
		t.Error("Expected:", true, "Got:", ok)
	}

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	// List accounts (expect 1)
	accounts, res, err = client.Accounts().List(context.Background())

	if http.StatusOK != res.StatusCode {
		t.Error("Expected:", http.StatusOK, "Got:", res.StatusCode)
	}

	if len(accounts) != 1 {
		t.Error("Expected: 1 account", "Got:", len(accounts))
	}

	// Attempt to fetch the deleted account
	account, res, err = client.Accounts().Fetch(context.Background(), uuid2)

	if http.StatusNotFound != res.StatusCode {
		t.Error("Expected:", http.StatusNotFound, "Got:", res.StatusCode)
	}

	// Note that this testing is non-exhaustive. Additionally we would seek to test:
	// - Failure paths
	// - Pagination
	// - All object attributes (e.g. https://github.com/google/go-github/blob/e5d8dd691c294eb6373a3dd5f58bec1e0d2ed3b1/github/reactions_test.go#L102-L105)
}

func generateAccount(uuid string) *form3.Account {
	return &form3.Account{
		ID:             uuid,
		OrganisationID: "358f775b-4ecd-4861-b33d-30df9a29de78",
		Type:           "accounts",
		Version:        0,
		Attributes: form3.AccountAttributes{
			Country: "GB",
			BankID:  "1112223",
		},
	}
}
