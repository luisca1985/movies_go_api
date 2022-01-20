CREATE DATABASE movies_db;

USE movies_db;

DROP TABLE IF EXISTS movie,
genre,
movie_genre;
CREATE TABLE movie (
  id INT AUTO_INCREMENT NOT NULL,
  title VARCHAR(128) NOT NULL,
  long_title VARCHAR(128) NOT NULL,
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
