# GO movies API
## Description

This API allows consumers to access the movies data. Most of the queries will be against local database but we use imdb npm package (https://godoc.org/github.com/eefret/gomdb) to query for movies when data isn’t available in the local database, which is stored for future reference. 

Movie object in database has the following properties as well as movie object in response of the apis:
- title
- long title
- released year
- rating
- id
- genres (array of strings)

### Update APIs
- Find movie by title by exact value that’s passed in the API. Notes: If there is no match in local database, imdb-api package is used for the search. If that returns result(s), then we store the result in database and return first value.
- Also we have an API that allows to update genres and ratings of the movie.

#### Search APIs 
These will be performed only against local database:
1. Search by Id
2. Search movies released in a particular year or given range
3. Search movies with rating higher or lower than passed in value.
4. Search movies with passed in genres value

## Prerequisites
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker compose](https://docs.docker.com/compose/install/)

## Get the project from Github
### Clone
```bash
git clone https://github.com/luisca1985/movies_go_api.git
```
### Access to the directory
```bash
cd movies_go_api.git
```

## Set Environment Variables
### OMDb API Key
#### Get OMDb API Key
Go to `http://omdbapi.com/apikey.aspx`, follow the instructions, generate a new API Key and activate this.

#### Include the API Key in environment variables
Open the file `app/.envs/.gomdb` with your editor and replace `<include_api_key>` with the key got previously.
```
# Golang Omdb API
# ------------------------------------------------------------------------------
# http://omdbapi.com/apikey.aspx
OMDB_API_KEY=<include_api_key>
```

### Mysql configuration
Open the file `sql/.envs/.mysql` with your editor and replace:
- `<include_mysql_root_password>`
- `<include_mysql_root_password>`
- `<include_mysql_user>`

Keep in mind that `MYSQL_ROOT_PASSWORD` and `MYSQL_PASSWORD` are different variables (although could have the same value).
```
# Mysql
# ------------------------------------------------------------------------------
MYSQL_ROOT_PASSWORD=<include_mysql_root_password>
MYSQL_PASSWORD=<include_mysql_password>
MYSQL_USER=<include_mysql_user>
MYSQL_DATABASE=movies_db
```

## Run the project with Docker
### Build, (re)create and start the container
```bash
docker-compose up --build -d
```
### Watch the logs
#### Go app
```bash
docker-compose logs -f app
```
#### Mysql database
```bash
docker-compose logs -f db
```
### Stop the project
```bash
docker-compose down
```

## API Endpoints

### Get a movie by `title`
Use  a `HTTP GET` request to the endpoint

```http
http:localhost/api/movies/title/<TITLE>
```

Gets a movie with the exact `<TITLE` in the mysql docker database. If the movie is not found, looks for it in `http://omdbapi.com/`, save the movie in the docker database and return it.

#### Example 1
```http
localhost/movies/title/back to the future
```
##### Response
```json
{
    "id": 1,
    "title": "back to the future",
    "long_title": "Back to the Future",
    "released_year": 1985,
    "rating": -1,
    "genres": [
        "Adventure",
        "Comedy",
        "Sci-Fi"
    ]
}
```
#### Example 2
```http
localhost/movies/title/terminator 2
```
##### Response
```json
{
    "id": 2,
    "title": "Terminator 2",
    "long_title": "Terminator 2: Judgment Day",
    "released_year": 1991,
    "rating": -1,
    "genres": [
        "Action",
        "Sci-Fi"
    ]
}
```

### Get a movie by `id`
Use  a `HTTP GET` request to the endpoint

```http
localhost/api/movies/<ID>
```

Gets a movie with the exact `<ID>` in the mysql docker database. 

#### Example
```http
localhost/movies/2
```
##### Response
```json
{
    "id": 2,
    "title": "terminator 2",
    "long_title": "Terminator 2: Judgment Day",
    "released_year": 1991,
    "rating": -1,
    "genres": [
        "Sci-Fi",
        "Action"
    ]
}
```

### List movies
Use  a `HTTP GET` request to the endpoint
```http
localhost/movies
```

Lists the movies in the mysql docker database.

#### Example
```http
http:localhost/movies
```
##### Response
```json
[
    {
        "id": 1,
        "title": "back to the future",
        "long_title": "Back to the Future",
        "released_year": 1985,
        "rating": -1,
        "genres": [
            "Adventure",
            "Comedy",
            "Sci-Fi"
        ]
    },
    {
        "id": 2,
        "title": "terminator 2",
        "long_title": "Terminator 2: Judgment Day",
        "released_year": 1991,
        "rating": -1,
        "genres": [
            "Sci-Fi",
            "Action"
        ]
    }
]
```

### Update a movie by `id`
Use  a `HTTP PUT` request to the endpoint

```http
localhost/movies/<id>
```
with content type `Content-Type: application/json`, and body
```json
{
	"rating": <RATING>,
	"genres": <[GENRES_LIST>
}
```

Update the movie `<RATING>` and `<[GENRES_LIST]>` list fields, using `<ID>` to select the movie. 

#### Example
```http
localhost/movies/1
```
body
```json
{
	"rating": 5,
	"genres": ["Sci=Fi", "Comedy"]
}
```
##### Response
```json
{
    "id": 1,
    "title": "back to the future",
    "long_title": "Back to the Future",
    "released_year": 1985,
    "rating": 5,
    "genres": [
        "Comedy",
        "Sci=Fi"
    ]
}
```
### Filters
Use  a `HTTP GET` request to the endpoints and the parameters to filter the movies.
#### Released between range
```http
localhost/movies?released_after=<MIN_DATE_VALUE>&released_before=<MAX_DATE_VALUE>
```
##### Example
```http
localhost/movies?released_after=1990&released_before=2000
```
###### Response
```json
[
    {
        "id": 2,
        "title": "terminator 2",
        "long_title": "Terminator 2: Judgment Day",
        "released_year": 1991,
        "rating": -1,
        "genres": [
            "Sci-Fi",
            "Action"
        ]
    },
    {
        "id": 5,
        "title": "jumanji",
        "long_title": "Jumanji",
        "released_year": 1995,
        "rating": -1,
        "genres": [
            "Adventure",
            "Comedy",
            "Family"
        ]
    }
]
```
#### Rating between range
```http
localhost/movies?rating_higher_than=<MIN_RATING_VALUE>&rating_lower_than=<MAX_RATING_VALUE>
```
##### Example
```http
localhost/movies?rating_higher_than=4
```
###### Response
```json
[
    {
        "id": 1,
        "title": "back to the future",
        "long_title": "Back to the Future",
        "released_year": 1985,
        "rating": 5,
        "genres": [
            "Comedy",
            "Sci=Fi"
        ]
    }
]
```
#### Genres
```http
localhost/movies?genre=<GENRE_VALUE>
```
##### Example
```http
localhost/movies?genre=Comedy
```
###### Response
```json
[
    {
        "id": 1,
        "title": "back to the future",
        "long_title": "Back to the Future",
        "released_year": 1985,
        "rating": 5,
        "genres": [
            "Comedy",
            "Sci=Fi"
        ]
    },
    {
        "id": 5,
        "title": "jumanji",
        "long_title": "Jumanji",
        "released_year": 1995,
        "rating": -1,
        "genres": [
            "Adventure",
            "Comedy",
            "Family"
        ]
    }
]
```