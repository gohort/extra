# extra

extra is used to flatten the extra fields in a JSON object that are not defined
in a Golang structure.

This is a functionality that `serde` in Rust provides called `flatten`, but is
not present in the default features of Golang JSON unmarshalling.

For example, if we have a JSON object for a person:

```json
{
    "firstName": "John",
    "lastName": "Doe",
    "age": 25,
    "favoriteColor": "red",
    "job": "mail man"
}
```

Now let's say we only want to capture John Doe's first and last name to use
in the code, and the rest we can ignore.

We want to ignore the other fields, but we don't want to omit them or delete them
because it's vital information about John Doe. So we create this structure:

```go
import "github.com/gohort/extra"

type Person struct {
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    // Any fields that are not included in the Go struct are placed in this map.
    extra.Any
}
```

So now when we marshal back using this structure it will be in the same format
as the JSON object that we received.

---

To see more examples please look in the `examples/` directory.