# Form3 Go API Client

### Usage:

```
import "github.com/benhawker/form3-api-client/form3"
```

Construct a new form3 client:
```
client, err := form3.NewClient(
    form3.SetScheme("https"),
    form3.SetHost("api.form3.tech"),
)
```

or just use the defaults (http://localhost:8080):

```
client, err := form3.NewClient()
```

Then use the `AccountsService` on the client to interact with `Account` resources.
```
// List all Accounts
accounts, _, err := client.Accounts().List(context.Background())
```


```
// List all Accounts with pagination (starting at page 2 (NB: zero indexed), with 30 results per page)
accounts, _, err := client.Accounts().Number(2).Size(30).List(context.Background())
```

```
// Fetch a single account by ID
account, res, err = client.Accounts().Fetch(context.Background(), "88cc4407-d170-44cd-b493-881edee7029c")
```

```
// Create a new account
acc := &form3.Account{
		ID:             "88cc4407-d170-44cd-b493-881edee7029c",
		OrganisationID: "88dd4407-d170-44cd-b493-881edee7029c",
		Type:           "accounts",
		Version:        0,
		Attributes: form3.AccountAttributes{
			Country: "GB",
			BankID: "1112223",
		},
	}

account, resp, err := client.Accounts().Create(context.Background(), newAcc)
```

```
// Delete a single account by ID
ok, resp, err := client.Accounts().Delete(context.Background(), "88cc4407-d170-44cd-b493-881edee7029c", 0)
```
NB: I believe there is an issue in the API. It always return `204 No Content` even when the UUID/version combination is not found).


### Testing:


From the root of this project run the integration tests:

```
$ docker-compose up
$ go test ./tests/... -v
```

Run the unit tests:
```
go test ./form3/... -v
```


### Suggested Improvements:
- Return a wrapped `*http.Response` rather than the raw `*http.Response`. The approach taken in the go-github client[here](https://github.com/google/go-github/blob/master/github/github.go#L404-L447) and [here](https://github.com/google/go-github/blob/master/github/github.go#L631-L656) and [here](https://github.com/google/go-github/blob/master/github/github.go#L768-L819) is one would consider.

- Allow easier extension of filtering (i.e. additional query string parameters beyong pagination). The current design can be improved.

- Use https://github.com/google/uuid for handling UUID's.

- Logging:
    - Use https://github.com/sirupsen/logrus or https://github.com/kataras/golog and use the log levels built into these libs.
    - Allow users of the library to configure logging verbosity
    - Remove logging by default from test output

- Allow configuration of timeouts.

- Taking inspiration from [this post](http://hassansin.github.io/Unit-Testing-http-client-in-Go) I decided to go with a unit testing approach using

- Use Table Driven tests. See [here](https://github.com/golang/go/wiki/TableDrivenTests) and [here](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go). An approach I would like to work into any refactor.

- Edge cases - there are **multiple unhandled cases** by both the Unit and Integration tests that I would hope to complete with more time. These would be essential to consider it production ready.


### Other notes:
- I have intentionally not handled "required attributes depending on the country the account is registered in" within the `Create` func, preferring to leave error handling to the server, and allow the client library to just relay any errors. No duplication of logic that can stray apart.