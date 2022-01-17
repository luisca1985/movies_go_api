package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"net/http"

	"github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

// movie represents data about a record album.
type Movie struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	ReleasedYear int      `json:"released_year"`
	Rating       float64  `json:"rating"`
	Genres       []string `json:"genres"`
}

var db *sql.DB

// albumsByArtist queries for albums that have the specified artist name.
func listMoviesByQuery(query string, args ...interface{}) ([]Movie, error) {
	// An albums slice to hold data from returned rows.

	var movies []Movie

	filterCriteria := "Criteria..."

	mov_rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("moviesCriteria %q: %v", filterCriteria, err)
	}
	defer mov_rows.Close()
	// Loop through mov_rows, using Scan to assign column data to struct fields.
	for mov_rows.Next() {
		var mov Movie
		if err := mov_rows.Scan(&mov.ID, &mov.Title, &mov.ReleasedYear, &mov.Rating); err != nil {
			return nil, fmt.Errorf("moviesCriteria %q: %v", filterCriteria, err)
		}

		gen_rows, err := db.Query("SELECT genre.genre FROM movie_genre JOIN genre ON movie_genre.genre_id = genre.id WHERE movie_genre.movie_id = ?", mov.ID)
		if err != nil {
			return nil, fmt.Errorf("movie %q: %v", mov.ID, err)
		}
		defer gen_rows.Close()

		for gen_rows.Next() {
			var genr string
			if err := gen_rows.Scan(&genr); err != nil {
				return nil, fmt.Errorf("Genre %q: %v", genr, err)
			}
			mov.Genres = append(mov.Genres, genr)
		}

		movies = append(movies, mov)
	}
	if err := mov_rows.Err(); err != nil {
		return nil, fmt.Errorf("moviesCriteria %q: %v", filterCriteria, err)
	}
	return movies, nil
}

func allMovies() ([]Movie, error) {
	movies, err := listMoviesByQuery("SELECT * FROM movie")
	return movies, err
}

func moviesByGenre(genre string) ([]Movie, error) {
	movies, err := listMoviesByQuery("SELECT movie.* FROM movie JOIN movie_genre ON movie.id = movie_genre.movie_id JOIN genre ON movie_genre.genre_id = genre.id WHERE genre.genre = ?", genre)
	return movies, err
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getMovies(c *gin.Context) {
	var movies []Movie
	var movies_filtered []Movie
	var err error

	id, id_err := strconv.Atoi(c.Param("id"))
	title := c.Query("title")
	genre := c.Query("genre")
	released_after, ra_err := strconv.Atoi(c.Query("released_after"))
	released_before, rb_err := strconv.Atoi(c.Query("released_before"))
	rating_higher_than, rh_err := strconv.ParseFloat(c.Query("rating_higher_than"), 64)
	rating_lower_than, rl_err := strconv.ParseFloat(c.Query("rating_lower_than"), 64)

	fmt.Println("Genre: ")
	fmt.Println(genre)
	if len(genre) > 0 {
		movies, err = moviesByGenre(genre)
	} else {
		movies, err = allMovies()
	}

	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Movies found: %v\n", movies)

	// c.IndentedJSON(http.StatusOK, movies)
	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, mov := range movies {
		if (mov.ID == id || id_err != nil) && (mov.Title == title || len(title) == 0) && (mov.ReleasedYear >= released_after || ra_err != nil) && (mov.ReleasedYear <= released_before || rb_err != nil) && (mov.Rating >= rating_higher_than || rh_err != nil) && (mov.Rating >= rating_lower_than || rl_err != nil) {
			fmt.Printf("Movie: %v\n", mov)
			movies_filtered = append(movies_filtered, mov)
		}
	}
	if len(movies_filtered) > 0 {
		c.IndentedJSON(http.StatusOK, movies_filtered)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "movies_db",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovies)
	router.Run("localhost:8080")
}
