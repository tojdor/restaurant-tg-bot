CREATE TABLE tables(
    number INTEGER NOT NULL UNIQUE,
    status TEXT NOT NULL CHECK(status IN('free', 'occupied'))
);