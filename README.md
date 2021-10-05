# Errors

A simple interface return descriptive errors from HTTP handlers. Expanding on the `error` interface, `Error` provides APIs to set the `type` and `message` of the error, as well as an optional `paramName` property.

## Get started

It's super simple to integrate with your APIs, first all install the module into your project.

```
> go get -u github.com/reecerussell/go-errors
```

Then as an example, here's how you can use it in your handlers.

```go

import (
    "net/http"

    "github.com/reecerussell/go-errors/errors"
)

...

func myHandler(w http.ResponseWriter, r *http.Request) {
    var data MyDataModel
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        // If the body is not valid, create a new validation error
        // with a type. Then write it to the response with the helper.
        err = errors.NewValidation(err).
            SetType("invalid body")

        errors.WriteResponse(w, err)
        return
    }

    err = MyProcess(&data) // returns standard error
    if err != nil {
        // The WriteResponse helper can also write standard errors.
        errors.WriteResponse(w, err)
    }
}

```
