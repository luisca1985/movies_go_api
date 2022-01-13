# Assignment

Develop APIs that allows consumers to access the movies data. Most of the queries will be against local database but you can use imdb npm package (https://godoc.org/github.com/eefret/gomdb ) to query for movies when data isn’t available in the local database and store it for future reference. Movie object in database will have following properties as well as movie object in response of the apis:

- title
- released year
- rating
- id
- genres (array of strings)

You will write API for following operations:

- Find movie by title by exact value that’s passed in the API. Notes: If there is no match in local database, use imdb-api package for the search. If that returns result(s), then store the result in database and return first value.
- API that allows updates to genres and ratings of the movie.

Implement following search APIs. These will be performed only against local database only:
1. Search by Id
2. Search movies released in a particular year or given range
3. Search movies with rating higher or lower than passed in value.
4. Search movies with passed in genres value

You can write 4 different APIs or write single API for all variations above.

You are free to make some decisions about the functionality on your own, there are no set
ground rules. After successfully completing the code execution, please upload the code on
github and send us the link. Also share postman collection with all the API urls so its easy to
validate. Thanks!