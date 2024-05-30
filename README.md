# mappath

mappath is a simple package for searching and modifying generic maps or slices using keypaths.

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
