# mappath

mappath is a simple package for searching and modifying generic maps (`map[string]any`) or slices (`[]any`) using keypaths.

This functionality was originally implemented in [Neptunus](https://github.com/gekatateam/neptunus) to navigate through event data.

## Examples
```go
rawJson := `
{
    "message": "user login",
    "metadata": {
        "user": {
            "name": "John Doe",
            "email": "johndoe@gmail.com",
            "roles": [ "employee", "manager" ],
            "age": 42
        }
    }
}
`

var data any
if err := json.Unmarshal([]byte(rawJson), &data); err != nil {
    return err
}

// get first user role
role, _ := mappath.Get(data, "metadata.user.roles.0")

// add login to user metadata
data, _ = mappath.Put(data, "metadata.user.login", "johndoe12")

// delete user second role
data, _ = mappath.Delete(data, "metadata.user.roles.1")
```

It may be not so conveniently, when function returns `nil` data if error occured. For this cases, use `Container`:

```go
rawJson := `
{
    "message": "user login",
    "metadata": {
        "user": {
            "name": "John Doe",
            "email": "johndoe@gmail.com",
            "roles": [ "employee", "manager" ],
            "age": 42
        }
    }
}
`

c := &mappath.Container{}
if err := json.Unmarshal([]byte(rawJson), &c.Data); err != nil {
    return err
}

// get first user role
role, err := c.Get("metadata.user.roles.0")

// add login to user metadata
err = c.Put("metadata.user.login", "johndoe12")

// delete user second role
err = c.Delete("metadata.user.roles.1")
```

`Container` stores data and updates it only if change operations have been performed successfully.
