CREATE TABLE IF NOT EXISTS actors (
    actor_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(1) CHECK (gender IN ('M', 'W')),
    birth_date DATE
);

CREATE TABLE IF NOT EXISTS movies (
    movie_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description VARCHAR(1000),
    release_date DATE,
    rating INT CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE movie_actors (
    movie_id BIGINT REFERENCES movies(movie_id),
    actor_id BIGINT REFERENCES actors(actor_id),
    PRIMARY KEY (movie_id, actor_id)
);