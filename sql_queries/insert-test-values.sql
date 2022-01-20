INSERT INTO
  movie (title, long_title, released_year, rating)
VALUES
  ('The Godfather​', 'The Godfather​',1972, 5.0),
  ('Kill Bill', 'Kill Bill',2004, 4.5),
  ('The Terminator', 'The Terminator',1984, 4.0),
  ('The Hangover', 'The Hangover', 2009, 3.5);
INSERT INTO
  genre (genre)
VALUES
  ('Drama​'),
  ('Crimen'),
  ('Acción'),
  ('Artes Marciales'),
  ('Ciencia Ficción'),
  ('Comedia');
INSERT INTO
  movie_genre (movie_id, genre_id)
VALUES
  (1, 1),
  (1, 2),
  (2, 3),
  (2, 4),
  (3, 3),
  (3, 5),
  (4, 6);