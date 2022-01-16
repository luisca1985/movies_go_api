package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

// movie represents data about a record album.
type Movie struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	ReleasedYear int      `json:"released_year"`
	Rating       float64  `json:"rating"`
	Genres       []string `json:"genres"`
}

var db *sql.DB

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

	//
	movies, err := moviesByGenre("Drama")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Movies found: %v\n", movies)
}

func moviesByGenre(genre string) ([]Movie, error) {
	movies, err := listMoviesQuery("SELECT movie.* FROM movie JOIN movie_genre ON movie.id = movie_genre.movie_id JOIN genre ON movie_genre.genre_id = genre.id WHERE genre.genre = ?", genre)
	return movies, err
}

// albumsByArtist queries for albums that have the specified artist name.
func listMoviesQuery(query string, args ...interface{}) ([]Movie, error) {
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
