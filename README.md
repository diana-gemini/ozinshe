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
4. http://localhost:3000/login
```json
{
    "email": "admin@mail.ru",
    "password": "Project2024&^!@"
}
```
5. http://localhost:3000/movie/create
```json
{
    "NameOfProject": "Superkolik samuryk",
	"Category":      "Telekhikaya, Multserial",
	"TypeOfProject": "Serial",
	"AgeCategory":   "8-10, 10-12",
	"Year":          "2020",
	"Timing":        "7",
	"Keywords":      "Tachka, Avto, Mult",
	"Description":   "Shytyrman ogigaly multserial...",
	"Director":      "Bakdaulet Alembekov",
	"Producer":      "Sandugash Kenzhebaeva",
    "CountOfSeason": {"1": ["vPQy7H-i3ww", "F_p7ePt17J4"], "2": ["dz8ET0_yzOM"]},
	"Cover":    "image link",
	"Screenshots":   ["image link 1","image link 2","image link 3"]
}
```
   
