## Ozinshe

### Steps 

Clone the repo

Rename the .env.example file to .env

Run the command `go run db/migrate/migrate.go` 

Run the project using the command `go run main.go`

### Routes
1. http://localhost:3000/signup
```json
{
    "email": "di@mail.ru",
    "password": "1478",
    "passwordrepeat": "1478"
}
```
2. http://localhost:3000/login
```json
{
    "email": "di@mail.ru",
    "password": "1478"
}
```
3. http://localhost:3000/logout

