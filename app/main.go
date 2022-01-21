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

	"github.com/eefret/gomdb"
)

var db *sql.DB

// Structs (Models) =========================

type Movie struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	LongTitle    string   `json:"long_title"`
	ReleasedYear int      `json:"released_year"`
	Rating       float64  `json:"rating"`
	Genres       []string `json:"genres"`
}

type Genre struct {
	ID    int    `json:"id"`
	Genre string `json:"genres"`
}

// Main Function =========================

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "db:3306",
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

// HTTP Functions =========================

func getMovies(c *gin.Context) {
	var movies []Movie
	movie := Movie{ID: -1}
	var movies_filtered []Movie
	var err error

	id, id_err := strconv.Atoi(c.Param("id"))
	title := c.Query("title")
	genre := c.Query("genre")
	released_after, ra_err := strconv.Atoi(c.Query("released_after"))
	released_before, rb_err := strconv.Atoi(c.Query("released_before"))
	rating_higher_than, rh_err := strconv.ParseFloat(c.Query("rating_higher_than"), 64)
	rating_lower_than, rl_err := strconv.ParseFloat(c.Query("rating_lower_than"), 64)

	if id_err == nil {
		var err_rmi error
		movie, err_rmi = retrieveMovieById(id)
		if err_rmi != nil {
			log.Fatal(fmt.Errorf("getMovies -> %v", err_rmi))
			return
		}
	} else if len(title) > 0 {
		var err_cmt error
		movie, err_cmt = createMovieWithTitleIfNotExist(title)
		if err_cmt != nil {
			log.Fatal(fmt.Errorf("getMovies -> %v", err_cmt))
		}
		if movie.ID != -1 {
			fmt.Printf("Movie %v get from http://omdbapi.com/", movie.Title)
		}
	} else if len(genre) > 0 {
		movies, err = listMoviesByGenre(genre)
		if err != nil {
			log.Fatal(fmt.Errorf("getMovies -> %v", err))
		}
	} else {
		movies, err = listAllMovies()
		if err != nil {
			log.Fatal(fmt.Errorf("getMovies -> %v", err))
		}
	}

	for _, mov := range movies {
		if (mov.ID == id || id_err != nil) && (mov.Title == title || len(title) == 0) && (mov.ReleasedYear >= released_after || ra_err != nil) && (mov.ReleasedYear <= released_before || rb_err != nil) && (mov.Rating >= rating_higher_than || rh_err != nil) && (mov.Rating >= rating_lower_than || rl_err != nil) {
			fmt.Printf("Movie: %v\n", mov)
			movies_filtered = append(movies_filtered, mov)
		}
	}

	if movie.ID != -1 {
		c.IndentedJSON(http.StatusOK, movie)
		return
	}
	if len(movies_filtered) > 0 {
		c.IndentedJSON(http.StatusOK, movies_filtered)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movies not found"})
}

func putMovies(c *gin.Context) {
	movieId, err_mid := strconv.Atoi(c.Param("id"))
	if err_mid != nil {
		fmt.Println(fmt.Errorf("putMovies -> %v", err_mid))
		return
	}

	var newMovie Movie
	newMovie.Rating = -1
	newMovie.Genres = []string{""}

	err_bind := c.BindJSON(&newMovie)
	if err_bind != nil {
		fmt.Println(fmt.Errorf("putMovies -> %v", err_bind))
		return
	}

	newRating := newMovie.Rating
	if newRating >= 0 {
		err_chr := updateMovieRatingByMovieId(movieId, newRating)
		if err_chr != nil {
			fmt.Println(fmt.Errorf("putMovies -> %v", err_chr))
			return
		}
	}

	newGenres := newMovie.Genres
	if len(newGenres) > 0 {
		err_chg := updateMovieGenresByMovieId(movieId, newGenres)
		if err_chg != nil {
			fmt.Println(fmt.Errorf("putMovies -> %v", err_chg))
			return
		}
	}

	movie, err_gmi := retrieveMovieById(movieId)
	if err_gmi != nil {
		fmt.Println(fmt.Errorf("putMovies -> %v", err_gmi))
		return
	}

	c.IndentedJSON(http.StatusCreated, movie)
}

// Update Functions =========================

func updateMovieRatingByMovieId(movieId int, newRating float64) error {
	query := `
	UPDATE movie 
	SET rating = ? 
	WHERE id = ?;`
	err := modifyTableByQuery(query, newRating, movieId)
	if err != nil {
		return fmt.Errorf("changeRatingMovieById -> %v", err)
	}
	return nil
}

func updateMovieGenresByMovieId(movieId int, newGenres []string) error {
	err_dlg := destroyGenresMovieById(movieId)
	if err_dlg != nil {
		return fmt.Errorf("changeGenresMovieById -> %v", err_dlg)
	}

	err_crg := createGenresIfNotExist(newGenres)
	if err_crg != nil {
		return fmt.Errorf("changeGenresMovieById -> %v", err_crg)
	}

	err_adg := createMovieGenresByMovieId(movieId, newGenres)
	if err_adg != nil {
		return fmt.Errorf("changeGenresMovieById -> %v", err_adg)
	}

	return nil
}

// Retrive Functions =========================

func retrieveMovieById(id int) (Movie, error) {
	query := `
	SELECT * 
	FROM movie 
	WHERE id = ?`
	movies, err := listMoviesByQuery(query, id)
	if err != nil {
		return Movie{ID: -1}, fmt.Errorf("movieById -> %v", err)
	}
	if len(movies) > 0 {
		return movies[0], nil
	}
	return Movie{ID: -1}, nil
}

func retrieveMovieByTitle(title string) (Movie, error) {
	query := `
	SELECT * 
	FROM movie 
	WHERE title = ? OR long_title = ?`
	movies, err := listMoviesByQuery(query, title, title)
	if err != nil {
		return Movie{ID: -1}, fmt.Errorf("retrieveMovieByTitle -> %v", err)
	}
	if len(movies) > 0 {
		return movies[0], nil
	}
	return Movie{ID: -1}, nil
}

func retriveMovieFromGomdbByTitle(title string) (Movie, error) {
	var gomdbMovie Movie
	api := gomdb.Init("4733c4c1")
	query := &gomdb.QueryData{Title: title}
	movieFromGomdb, err := api.MovieByTitle(query)
	if err != nil {
		fmt.Println(fmt.Errorf("Error retrieving movie from Golang Omdb API -> %v", err))
		return Movie{ID: -1}, nil
	}

	var err_ryr error
	gomdbMovie.Title = title
	gomdbMovie.LongTitle = movieFromGomdb.Title
	gomdbMovie.Rating = -1
	gomdbMovie.Genres = strings.Split(strings.ReplaceAll(movieFromGomdb.Genre, " ", ""), ",")
	gomdbMovie.ReleasedYear, err_ryr = strconv.Atoi(movieFromGomdb.Year)
	if err_ryr != nil {
		return Movie{ID: -1}, fmt.Errorf("retriveMovieFromGomdbByTitle -> %v", err_ryr)
	}
	return gomdbMovie, nil
}

// Destroy Functions =========================

func destroyGenresMovieById(movieId int) error {
	query := `
	DELETE FROM movie_genre 
	WHERE movie_id = ?;`
	err := modifyTableByQuery(query, movieId)
	if err != nil {
		return fmt.Errorf("deleteGenresMovieById: movie id %q -> %v", movieId, err)
	}
	return nil
}

// Create Functions =========================

func createMovieWithTitleIfNotExist(title string) (Movie, error) {
	localMovie, err_loc := retrieveMovieByTitle(title)
	if err_loc != nil {
		return Movie{ID: -1}, fmt.Errorf("createMovieWithTitleIfNotExist: movie title %q -> %v", title, err_loc)
	}

	if len(localMovie.Title) == 0 {
		gomdbMovie, err_rem := retriveMovieFromGomdbByTitle(title)
		if err_rem != nil {
			return Movie{ID: -1}, fmt.Errorf("createMovieWithTitleIfNotExist: movie title %q -> %v", title, err_rem)
		}
		if gomdbMovie.ID != -1 {
			err_crm := createMovie(gomdbMovie)
			if err_crm != nil {
				return Movie{ID: -1}, fmt.Errorf("createMovieWithTitleIfNotExist: movie title %q -> %v", gomdbMovie.Title, err_crm)
			}

			newMovieWithoutGenres, err_new := retrieveMovieByTitle(gomdbMovie.Title)
			if err_new != nil {
				return Movie{ID: -1}, fmt.Errorf("createMovieWithTitleIfNotExist: new movie title %q -> %v", title, err_new)
			}

			err_chg := updateMovieGenresByMovieId(newMovieWithoutGenres.ID, gomdbMovie.Genres)
			if err_chg != nil {
				return Movie{ID: -1}, fmt.Errorf("createMovieWithTitleIfNotExist -> %v", err_chg)
			}

			newMovieWithGenres, err_new := retrieveMovieById(newMovieWithoutGenres.ID)
			if err_new != nil {
				return Movie{ID: -1}, fmt.Errorf("createMovieWithTitleIfNotExist: new movie title %q -> %v", title, err_new)
			}

			return newMovieWithGenres, nil
		}
		return Movie{ID: -1}, nil
	}

	return localMovie, nil
}

func createMovie(movie Movie) error {
	query := `
	INSERT INTO movie (title, long_title, released_year, rating)
	VALUES (?,?,?,?)`
	err := modifyTableByQuery(query, movie.Title, movie.LongTitle, movie.ReleasedYear, movie.Rating)
	if err != nil {
		fmt.Println(fmt.Errorf("createMovie: title %q -> %v", movie.Title, err))
	}

	return nil
}

func createMovieGenresByMovieId(movieId int, newGenresStr []string) error {
	Genres, err_fgen := listGenresInStringList(newGenresStr)
	if err_fgen != nil {
		return fmt.Errorf("addNewGenresToMovieById -> %v", err_fgen)
	}
	for _, genre := range Genres {
		err_ing := createGenreInMovie(movieId, genre.ID)
		if err_ing != nil {
			return fmt.Errorf("addNewGenresToMovieById -> %v", err_ing)
		}
	}
	return nil
}

func createGenreInMovie(movieId int, genreId int) error {
	query := `
	INSERT INTO movie_genre (movie_id, genre_id)
	VALUES (?,?)`
	err := modifyTableByQuery(query, movieId, genreId)
	if err != nil {
		return fmt.Errorf("insertGenreInMovie:movie %q, genre %q -> %v", movieId, genreId, err)
	}
	return nil
}

func createGenresIfNotExist(newGenres []string) error {
	notExistGenres, err_neg := listNotExistGenresInStrinList(newGenres)
	if err_neg != nil {
		return fmt.Errorf("createGenresIfNotExist -> %v", err_neg)
	}

	for _, genre := range notExistGenres {
		createGenre(genre)
		if err_neg != nil {
			return fmt.Errorf("createGenresIfNotExist -> %v", err_neg)
		}
	}
	return nil
}

func createGenre(genre string) error {
	query := `
	INSERT INTO genre (genre)
	VALUES (?)`
	err := modifyTableByQuery(query, genre)
	if err != nil {
		return fmt.Errorf("insertGenre: genre %q -> %v", genre, err)
	}

	return nil
}

// List Functions =========================

func listAllMovies() ([]Movie, error) {
	query := `
	SELECT * 
	FROM movie`
	movies, err := listMoviesByQuery(query)
	if err != nil {
		return nil, fmt.Errorf("allMovies -> %v", err)
	}
	return movies, nil
}

func listMoviesByGenre(genre string) ([]Movie, error) {
	query := `
	SELECT movie.* 
	FROM movie 
	JOIN movie_genre ON movie.id = movie_genre.movie_id 
	JOIN genre ON movie_genre.genre_id = genre.id 
	WHERE genre.genre = ?`
	movies, err := listMoviesByQuery(query, genre)
	if err != nil {
		return movies, fmt.Errorf("moviesByGenre -> %v", err)
	}
	return movies, nil
}

func listGenresByMovieId(id int) ([]Genre, error) {
	query := `
	SELECT genre.* 
	FROM movie_genre 
	JOIN genre ON movie_genre.genre_id = genre.id 
	WHERE movie_genre.movie_id = ?`
	genres, err := listGenresByQuery(query, id)
	if err != nil {
		return nil, fmt.Errorf("genresByMovieId -> %v", err)
	}
	return genres, nil
}

func listAllGenres() ([]Genre, error) {
	query := `
	SELECT * 
	FROM genre`
	genres, err := listGenresByQuery(query)
	if err != nil {
		return nil, fmt.Errorf("allGenres -> %v", err)
	}
	return genres, nil
}

func listGenresInStringList(genres_str []string) ([]Genre, error) {
	query := `
	SELECT * 
	FROM genre 
	WHERE genre IN (?` + strings.Repeat(`,?`, len(genres_str)-1) + `)`
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

func listNotExistGenresInStrinList(newGenres_str []string) ([]string, error) {
	var notExistGenres []string
	existGenres, err := listAllGenres()
	if err != nil {
		return nil, fmt.Errorf("getNotExistGenres -> %v", err)
	}

	for _, newGenre := range newGenres_str {
		exist := false
		for _, existGenre := range existGenres {
			if newGenre == existGenre.Genre {
				exist = true
				break
			}
		}
		if !exist {
			notExistGenres = append(notExistGenres, newGenre)
		}
	}

	return notExistGenres, nil
}

// Low Level Query Functions =========================

func modifyTableByQuery(query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("updateAndDeleteRecords -> %v", err)
	}
	return nil
}

func listMoviesByQuery(query string, args ...interface{}) ([]Movie, error) {

	var movies []Movie

	mov_rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listMoviesByQuery: Movies Query -> %v", err)
	}
	defer mov_rows.Close()

	for mov_rows.Next() {
		var mov Movie
		err_scan := mov_rows.Scan(&mov.ID, &mov.Title, &mov.LongTitle, &mov.ReleasedYear, &mov.Rating)
		if err_scan != nil {
			return nil, fmt.Errorf("listMoviesByQuery: Movie %q Scan -> %v", mov.Title, err_scan)
		}

		genres, err_gen := listGenresByMovieId(mov.ID)
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
