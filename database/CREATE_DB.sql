USE movieratings;

DROP TABLE IF EXISTS MovieCredits;
DROP TABLE IF EXISTS SeriesCredits;
DROP TABLE IF EXISTS MovieWeekPerson;
DROP TABLE IF EXISTS PersonWeek;
DROP TABLE IF EXISTS Personen;
DROP TABLE IF EXISTS SeriesWeekPopularity;
DROP TABLE IF EXISTS MovieWeekPopularity;
DROP TABLE IF EXISTS SeriesGenre;
DROP TABLE IF EXISTS MovieGenre;
DROP TABLE IF EXISTS MovieCountry;
DROP TABLE IF EXISTS SeriesCountry;
DROP TABLE IF EXISTS SeriesNetwork;
DROP TABLE IF EXISTS Genres;
DROP TABLE IF EXISTS Series;
DROP TABLE IF EXISTS Movies;
DROP TABLE IF EXISTS Networks;
DROP TABLE IF EXISTS Countries;

CREATE TABLE Countries (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    cname VARCHAR(35)
);

CREATE TABLE Networks (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    nname VARCHAR(50),
    logo VARCHAR(50),
    originCountry VARCHAR(10),
    FOREIGN KEY (originCountry) REFERENCES Countries(id)
);

CREATE TABLE Movies (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    title VARCHAR(100),
    overview VARCHAR(1000),
    popularity DOUBLE,
    revenue VARCHAR(50),
    posterPath VARCHAR(50),
    releaseDate VARCHAR(10),
    voteAvg DOUBLE,
    voteCount INT,
    runtime INT,
    tagline VARCHAR(75)
);

CREATE TABLE Genres (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    genre VARCHAR(50)
);

CREATE TABLE MovieGenre (
	movieId VARCHAR(10) NOT NULL,
    genreId VARCHAR(10) NOT NULL,
    PRIMARY KEY (movieId, genreId),
    FOREIGN KEY (movieId) REFERENCES Movies(id),
    FOREIGN KEY (genreId) REFERENCES Genres(id)
);

CREATE TABLE MovieCountry (
	movieId VARCHAR(10) NOT NULL,
    countryId VARCHAR(10) NOT NULL,
    PRIMARY KEY (movieId, countryId),
    FOREIGN KEY (movieId) REFERENCES Movies(id),
    FOREIGN KEY (countryId) REFERENCES Countries(id)
);

CREATE TABLE MovieWeekPopularity (
	movieId VARCHAR(10) NOT NULL,
    weekNr VARCHAR(10) NOt NULL,
    popularity DOUBLE,
    voteAvg DOUBLE,
    voteCount INT,
    FOREIGN KEY (movieId) REFERENCES Movies(id),
    PRIMARY KEY (movieId, weekNr)
);

CREATE TABLE Series (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    title VARCHAR(100),
    overview VARCHAR(1000),
    popularity DOUBLE,
    seasons INT,
    episodes INT,
    posterPath VARCHAR(50),
    voteAvg DOUBLE,
    voteCount INT,
    firstAir VARCHAR(10),
    lastAir VARCHAR(10),
    tagline VARCHAR(100)
);

CREATE TABLE SeriesGenre (
	seriesId VARCHAR(10) NOT NULL,
    genreId VARCHAR(10) NOT NULL,
    PRIMARY KEY (seriesId, genreId),
    FOREIGN KEY (seriesId) REFERENCES Series(id),
    FOREIGN KEY (genreId) REFERENCES Genres(id)
);

CREATE TABLE SeriesCountry (
	seriesId VARCHAR(10) NOT NULL,
    countryId VARCHAR(10) NOT NULL,
    PRIMARY KEY (seriesId, countryId),
    FOREIGN KEY (seriesId) REFERENCES Series(id),
    FOREIGN KEY (countryId) REFERENCES Countries(id)
);

CREATE TABLE SeriesNetwork (
	seriesId VARCHAR(10) NOT NULL,
    networkId VARCHAR(10) NOT NULL,
    PRIMARY KEY (seriesId, networkId),
    FOREIGN KEY (seriesId) REFERENCES Series(id),
    FOREIGN KEY (networkId) REFERENCES Networks(id)
);

CREATE TABLE SeriesWeekPopularity (
	seriesId VARCHAR(10) NOT NULL,
    weekNr VARCHAR(10) NOT NULL,
    popularity DOUBLE,
    voteAvg DOUBLE,
    voteCount INT,
    FOREIGN KEY (seriesId) REFERENCES Series(id),
    PRIMARY KEY (seriesId, weekNr)
);

CREATE TABLE Personen (
	id VARCHAR(10) PRIMARY KEY NOT NULL,
    name VARCHAR(50),
    birthday VARCHAR(10),
    deathday VARCHAR(10),
    popularity DOUBLE,
    profilePath VARCHAR(50),
    gender INT,
    profession VARCHAR(25)
);

CREATE TABLE PersonWeek (
	personId VARCHAR(10) NOT NULL,
    weekNr int NOT NULL,
    PRIMARY KEY (personId, weekNr),
    FOREIGN KEY (personID) REFERENCES Personen(id)
);

CREATE TABLE MovieWeekPerson (
	movieId VARCHAR(10) NOT NULL,
    weekNr VARCHAR(10) NOt NULL,
    personId VARCHAR(10) NOT NULL,
    popularity DOUBLE,
    revenue DOUBLE,
    voteAvg DOUBLE,
    voteCount INT,
    FOREIGN KEY (movieId) REFERENCES Movies(id),
    FOREIGN KEY (personId) REFERENCES Personen(id),
    PRIMARY KEY (movieId, weekNr, personId)
);

CREATE TABLE MovieCredits (
	movieId VARCHAR(10) NOT NULL,
    personId VARCHAR(10) NOT NULL,
    job varchar(120) NOT NULL,
    PRIMARY KEY (movieId, personId, job),
    FOREIGN KEY (movieId) REFERENCES Movies(id),
    FOREIGN KEY (personId) REFERENCES Personen(id)
);

CREATE TABLE SeriesCredits (
	seriesId VARCHAR(10) NOT NULL,
    personId VARCHAR(10) NOT NULL,
    job varchar(120) NOT NULL,
    PRIMARY KEY (seriesId, personId, job),
    FOREIGN KEY (seriesId) REFERENCES Series(id),
    FOREIGN KEY (personId) REFERENCES Personen(id)
);