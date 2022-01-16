SELECT
  movie.*
--   GROUP_CONCAT(genre.genre) AS genres
FROM
  movie
  JOIN movie_genre ON movie.id = movie_genre.movie_id
  JOIN genre ON movie_genre.genre_id = genre.id
WHERE
  genre.genre = 'Acción'
-- GROUP BY
--   movie.title,
--   movie.released_year,
--   movie.rating;