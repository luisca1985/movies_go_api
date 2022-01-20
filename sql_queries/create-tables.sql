DROP TABLE IF EXISTS movie,
genre,
movie_genre;
CREATE TABLE movie (
  id INT AUTO_INCREMENT NOT NULL,
  title VARCHAR(128) NOT NULL,
  released_year INT NOT NULL,
  rating DECIMAL(2, 1) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (title)
);
CREATE TABLE genre (
  id INT AUTO_INCREMENT NOT NULL,
  genre VARCHAR(128) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (genre)
);
CREATE TABLE movie_genre (
  -- id INT AUTO_INCREMENT NOT NULL,
  movie_id INT NOT NULL,
  genre_id INT NOT NULL,
  -- PRIMARY KEY (`id`),
  PRIMARY KEY (movie_id,genre_id),
  CONSTRAINT moviegenre_movie FOREIGN KEY (movie_id) REFERENCES movie (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT moviegenre_genre FOREIGN KEY (genre_id) REFERENCES genre (id) ON DELETE CASCADE ON UPDATE CASCADE
);
ALTER TABLE
  movie_genre
ADD
  (
    INDEX moviegenre_movie_idx (movie_id ASC),
    INDEX moviegenre_genre_idx (genre_id ASC)
  );
INSERT INTO
  movie (title, released_year, rating)
VALUES
  ('The Godfather​', 1972, 5.0),
  ('Kill Bill', 2004, 4.5),
  ('The Terminator', 1984, 4.0),
  ('The Hangover', 2009, 3.5);
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