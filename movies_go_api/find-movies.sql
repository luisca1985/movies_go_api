SELECT movie.title, movie.released_year, movie.rating, GROUP_CONCAT(genre.genre)
FROM movie 
JOIN movie_genre ON movie.id = movie_genre.movie_id
JOIN genre ON movie_genre.genre_id = genre.id
GROUP BY movie.title, movie.released_year, movie.rating;