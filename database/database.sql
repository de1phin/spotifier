DROP DATABASE IF EXISTS Spotifier;

CREATE DATABASE Spotifier;
use Spotifier;

CREATE TABLE Users (
    id INT PRIMARY KEY,
    spotify_id VARCHAR(256),
    name VARCHAR(256),
    access_token VARCHAR(256),
    refresh_token VARCHAR(256),
    expires_after DATETIME,
    last_update BIGINT
);

CREATE TABLE TrackStatistics (
    user_id INT,
    track_id VARCHAR(256),
    time_listened BIGINT,
    statistic_date DATE
);

CREATE TABLE ArtistStatistics (
    user_id INT,
    artist_id VARCHAR(256),
    time_listened BIGINT,
    statistic_date DATE
);

CREATE TABLE AlbumStatistics (
    user_id INT,
    album_id VARCHAR(256),
    time_listened BIGINT,
    statistic_date DATE
);