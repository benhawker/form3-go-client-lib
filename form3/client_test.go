package form3

import (
	"context"
	"testing"
)

func TestClientDefaults(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client.scheme != defaultScheme {
		t.Errorf("expected scheme to be %s got: %q", defaultScheme, client.scheme)
	}
	if client.host != defaultHost {
		t.Errorf("expected host to be %s; got: %s", defaultHost, client.host)
	}
}

func TestMakeRequest(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.MakeRequest(context.TODO(), MakeRequestOptions{
		Method: "GET",
		Path:   "/organisation/accounts",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response to be != nil")
	}
}
