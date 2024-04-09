## Ozinshe

### Steps 

Clone the repo

Rename the .env.example file to .env

Run the command `go run db/migrate/migrate.go` 

Run the project using the command `go run main.go`

Test routes in Postman (my collection in Postman - https://speeding-flare-870218.postman.co/workspace/Go~fe31b55f-1eb9-480b-837e-1a4dcfaea319/collection/22575040-46b5534d-d00d-4d00-b864-c7001da97514?action=share&creator=22575040)

### Routes
1. http://localhost:3000/signup - Signup user
```json
{
    "email": "di@mail.ru",
    "password": "1478",
    "passwordrepeat": "1478"
}
```
2. http://localhost:3000/login - Login user
```json
{
    "email": "di@mail.ru",
    "password": "1478"
}
```
3. http://localhost:3000/logout - Logout
4. http://localhost:3000/login - Login Admin
```json
{
    "email": "admin@mail.ru",
    "password": "Project2024&^!@"
}
```
5. http://localhost:3000/movie/create - Create movie only for Admin
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
        "CountOfSeason": [
    {
      "season": "1",
      "linkOfSeries": ["vPQy7H-i3ww", "F_p7ePt17J4"]
    },
    {
      "season": "2",
      "linkOfSeries": ["dz8ET0_yzOM"]
    }
  ],
        "Cover":    "image link",
        "Screenshots":   ["image link 1","image link 2","image link 3"]
    }
```
6. http://localhost:3000/movie/1/edit - Edit movie only for Admin
7. http://localhost:3000/movie/1/update - Update movie only for Admin
```json
    {
        "NameOfProject": "Superkolik samuryk new",
        "Category":      "Telekhikaya, Multserial",
        "TypeOfProject": "Serial",
        "AgeCategory":   "8-10, 10-12",
        "Year":          "2020",
        "Timing":        "7",
        "Keywords":      "Tachka, Avto, Mult",
        "Description":   "Shytyrman ogigaly multserial...",
        "Director":      "Bakdaulet Alembekov",
        "Producer":      "Sandugash Kenzhebaeva",
        "CountOfSeason": [
    {
      "season": "1",
      "linkOfSeries": ["vPQy7H-i3ww", "F_p7ePt17J4"]
    },
    {
      "season": "2",
      "linkOfSeries": ["dz8ET0_yzOM"]
    }
  ],
        "Cover":    "image link",
        "Screenshots":   ["image link 1","image link 2","image link 3"]
    }
```
8. http://localhost:3000/movie/1/delete - Delete movie only for Admin
9. http://localhost:3000/all - Get all movies
10. http://localhost:3000/movie/1 - Get movie by ID
11. http://localhost:3000/movie/1/series/1/2 - Get series by ID
12. http://localhost:3000/trends - Get trends
13. http://localhost:3000/newprojects - Get new projects
14. http://localhost:3000/telehikaya - Get telehikaya
   
