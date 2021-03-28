USE movieratings;

DROP TABLE IF EXISTS film_week_person;
DROP TABLE IF EXISTS person;
DROP TABLE IF EXISTS series_week_popularity;
DROP TABLE IF EXISTS film_Week_popularity;
DROP TABLE IF EXISTS series_Genre;
DROP TABLE IF EXISTS film_genre;
DROP TABLE IF EXISTS genre;
DROP TABLE IF EXISTS series;
DROP TABLE IF EXISTS film;

CREATE TABLE film (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    title VARCHAR(50),
    overview VARCHAR(50),
    popularity DOUBLE,
    revenue DOUBLE,
    posterPath VARCHAR(50),
    release_date DATE,
    voteAvg DOUBLE,
    voteCount INT,
    inProduction BOOLEAN
);

CREATE TABLE genre (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    genre VARCHAR(50)
);

CREATE TABLE film_genre (
	filmId VARCHAR(10) NOT NULL,
    genreId VARCHAR(10) NOT NULL,
    PRIMARY KEY (filmId, genreId),
    FOREIGN KEY (filmId) REFERENCES film(id),
    FOREIGN KEY (genreId) REFERENCES genre(id)
);

CREATE TABLE film_week_popularity (
	filmId VARCHAR(10) NOT NULL,
    weekNr INT NOt NULL,
    popularity DOUBLE,
    revenue DOUBLE,
    voteAvg DOUBLE,
    voteCount INT,
    FOREIGN KEY (filmId) REFERENCES film(id),
    PRIMARY KEY (filmId, weekNr)
);

CREATE TABLE series (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    title VARCHAR(50),
    overview VARCHAR(50),
    popularity DOUBLE,
    seasons INT,
    episodes INT,
    posterPath VARCHAR(50),
    release_date DATE,
    voteAvg DOUBLE,
    voteCount INT,
    inProduction BOOLEAN
);

CREATE TABLE series_genre (
	seriesId VARCHAR(10) NOT NULL,
    genreId VARCHAR(10) NOT NULL,
    PRIMARY KEY (seriesId, genreId),
    FOREIGN KEY (seriesId) REFERENCES series(id),
    FOREIGN KEY (genreId) REFERENCES genre(id)
);

CREATE TABLE series_week_popularity (
	seriesId VARCHAR(10) NOT NULL,
    weekNr INT NOT NULL,
    popularity DOUBLE,
    voteAvg DOUBLE,
    voteCount INT,
    FOREIGN KEY (seriesId) REFERENCES series(id),
    PRIMARY KEY (seriesId, weekNr)
);

CREATE TABLE person (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    name VARCHAR(50),
    birthday DATE,
    deathday DATE,
    popularity DOUBLE,
    profilePath VARCHAR(50),
    gender INT,
    profession VARCHAR(25)
);

CREATE TABLE Film_Week_Person (
	filmId VARCHAR(10) NOT NULL,
    weekNr INT NOt NULL,
    personId VARCHAR(10) NOT NULL,
    popularity DOUBLE,
    revenue DOUBLE,
    voteAvg DOUBLE,
    voteCount INT,
    FOREIGN KEY (filmId) REFERENCES film(id),
    FOREIGN KEY (personId) REFERENCES person(id),
    PRIMARY KEY (filmId, weekNr, personId)
);