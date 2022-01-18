package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

type Genre struct {
	ID    int    `json:"id"`
	Genre string `json:"genres"`
}

var db *sql.DB

func listGenresByQuery(query string, args ...interface{}) ([]Genre, error) {
	var genres []Genre
	gen_rows, err_query := db.Query(query, args...)
	if err_query != nil {
		return nil, fmt.Errorf("listGenresByQuery: Genres Query -> %w", err_query)
	}
	defer gen_rows.Close()

	for gen_rows.Next() {
		var genre Genre
		err_scan := gen_rows.Scan(&genre.ID, &genre.Genre)
		if err_scan != nil {
			return nil, fmt.Errorf("listGenresByQuery: Genre %q Scan -> %w", genre.Genre, err_scan)
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func genresByMovieId(id int) ([]Genre, error) {
	genres, err := listGenresByQuery("SELECT genre.* FROM movie_genre JOIN genre ON movie_genre.genre_id = genre.id WHERE movie_genre.movie_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("genresByMovieId -> %v", err)
	}
	return genres, nil
}

func allGenres() ([]Genre, error) {
	genres, err := listGenresByQuery("SELECT * FROM genre")
	if err != nil {
		return nil, fmt.Errorf("allGenres -> %v", err)
	}
	return genres, nil
}

func findGenres(genres_str []string) ([]Genre, error) {
	query := "SELECT * FROM genre WHERE genre IN (?" + strings.Repeat(",?", len(genres_str)-1) + ")"
	genres_inter := make([]interface{}, len(genres_str))
	for i, v := range genres_str {
		genres_inter[i] = v
	}
	genres_in_db, err := listGenresByQuery(query, genres_inter...)
	fmt.Println("Genres: ")
	for _, genre := range genres_in_db {
		fmt.Println(genre.Genre)
	}
	if err != nil {
		return nil, fmt.Errorf("findGenres -> %v", err)
	}
	return genres_in_db, err
}

// albumsByArtist queries for albums that have the specified artist name.
func listMoviesByQuery(query string, args ...interface{}) ([]Movie, error) {
	// An albums slice to hold data from returned rows.

	var movies []Movie

	mov_rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listMoviesByQuery: Movies Query -> %v", err)
	}
	defer mov_rows.Close()
	// Loop through mov_rows, using Scan to assign column data to struct fields.
	for mov_rows.Next() {
		var mov Movie
		err_scan := mov_rows.Scan(&mov.ID, &mov.Title, &mov.ReleasedYear, &mov.Rating)
		if err_scan != nil {
			return nil, fmt.Errorf("listMoviesByQuery: Movie %q Scan -> %v", mov.Title, err_scan)
		}

		genres, err_gen := genresByMovieId(mov.ID)
		for _, genre := range genres {
			mov.Genres = append(mov.Genres, genre.Genre)
		}
		if err_gen != nil {
			return nil, fmt.Errorf("listMoviesByQuery -> %v", err_gen)
		}

		movies = append(movies, mov)
	}
	err_err := mov_rows.Err()
	if err != nil {
		return nil, fmt.Errorf("listMoviesByQuery: Err -> %v", err_err)
	}
	return movies, nil
}

func allMovies() ([]Movie, error) {
	movies, err := listMoviesByQuery("SELECT * FROM movie")
	if err != nil {
		return nil, fmt.Errorf("allMovies -> %v", err)
	}
	return movies, nil
}

func movieById(id int) (Movie, error) {
	movies, err := listMoviesByQuery("SELECT * FROM movie WHERE id = ?", id)
	if err != nil {
		return Movie{}, fmt.Errorf("movieById -> %v", err)
	}
	return movies[0], nil
}

func moviesByGenre(genre string) ([]Movie, error) {
	movies, err := listMoviesByQuery("SELECT movie.* FROM movie JOIN movie_genre ON movie.id = movie_genre.movie_id JOIN genre ON movie_genre.genre_id = genre.id WHERE genre.genre = ?", genre)
	if err != nil {
		return movies, fmt.Errorf("moviesByGenre -> %v", err)
	}
	return movies, nil
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
		if err != nil {
			log.Fatal(fmt.Errorf("getMovies -> %v", err))
		}
	} else {
		movies, err = allMovies()
		if err != nil {
			log.Fatal(fmt.Errorf("getMovies -> %v", err))
		}
	}
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

func updateAndDeleteRecords(query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("updateAndDeleteRecords -> %v", err)
	}
	return nil
}

func changeRatingMovieById(movieId int, newRating float64) error {
	query := `
	UPDATE movie 
	SET rating = ? 
	WHERE id = ?;`
	err := updateAndDeleteRecords(query, newRating, movieId)
	if err != nil {
		return fmt.Errorf("changeRatingMovieById -> %v", err)
	}
	return nil
}

func deleteGenresMovieById(movieId int) error {
	query := `DELETE movie SET rating = ? WHERE id = ?;`
	err := updateAndDeleteRecords(query, newRating, movieId)
	if err != nil {
		return fmt.Errorf("changeRatingMovieById -> %v", err)
	}
	return nil
}

func changeGenresMovieById(movieId int, newGenres []string) error {
	err_dlg := deleteGenresMovieById(movieId)
	if err_dlg != nil {
		return fmt.Errorf("changeRatingMovieById -> %v", err_dlg)
	}

	return nil
}

// func removeAllGenresByMovieId(id int) error {

// 	return nil
// }
// postAlbums adds an album from JSON received in the request body.
func putMovies(c *gin.Context) {
	movieId, err_mid := strconv.Atoi(c.Param("id"))
	if err_mid != nil {
		fmt.Println(fmt.Errorf("putMovies -> %v", err_mid))
		return
	}

	var newMovie Movie
	newMovie.Rating = -1
	newMovie.Genres = []string{""}

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	err_bind := c.BindJSON(&newMovie)
	if err_bind != nil {
		fmt.Println(fmt.Errorf("putMovies -> %v", err_bind))
		return
	}

	newRating := newMovie.Rating
	if newRating >= 0 {
		err_chr := changeRatingMovieById(movieId, newRating)
		if err_chr != nil {
			fmt.Println(fmt.Errorf("putMovies -> %v", err_chr))
			return
		}
	}

	newGenres := newMovie.Genres
	if len(newGenres) > 0 {
		err_chg := changeGenresMovieById(movieId, newGenres)
		if err_chg != nil {
			fmt.Println(fmt.Errorf("putMovies -> %v", err_chg))
			return
		}
	}

	// movieId, err_mid := strconv.Atoi(c.Param("id"))
	// if err_mid != nil {
	// 	fmt.Println(fmt.Errorf("putMovies -> %v", err_mid))
	// 	return
	// }

	// oldMovie, err_old := movieById(movieId)
	// if err_old != nil {
	// 	fmt.Println(fmt.Errorf("putMovies -> %v", err_old))
	// 	return
	// }

	// if len(oldMovie) == 0 {
	// 	return
	// }

	// oldGenres, err_oldg := genresByMovieId(movieId)
	// if err_oldg != nil {
	// 	fmt.Println(fmt.Errorf("putMovies -> %v", err_oldg))
	// 	return
	// }

	// for _, oldGenre := range oldGenres {
	// 	oldInNew := false
	// 	for _, newGenre := range newMovie.Genres {
	// 		if oldGenre.Genre == newGenre {
	// 			oldInNew = true
	// 			break
	// 		}
	// 	}
	// 	if !oldInNew {
	// 		// Se borra el genero viejo
	// 	}
	// }

	// newGenres, err_gen := findGenres(newMovie.Genres)
	// if err_gen != nil {
	// 	fmt.Println(fmt.Errorf("putMovies -> %v", err_gen))
	// }

	// for _, newGenre := range newMovie.Genres {
	// }

	// fmt.Println(newMovie)
	// fmt.Println(genres)

	// Add the new album to the slice.
	// albums = append(albums, newAlbum)
	// c.IndentedJSON(http.StatusCreated, genres)
	movie, err_gmi := movieById(movieId)
	if err_bind != nil {
		fmt.Println(fmt.Errorf("putMovies -> %v", err_gmi))
		return
	}

	c.IndentedJSON(http.StatusCreated, movie)
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
	router.PUT("/movies/:id", putMovies)
	router.Run("localhost:8080")
}
