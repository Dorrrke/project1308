CREATE TABLE IF NOT EXISTS users (
    uid varchar(36) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    age int NOT NULL,
    phone text,
    UNIQUE (email),
    UNIQUE (phone)
);

CREATE TABLE IF NOT EXISTS cars (
    cid varchar(36) NOT NULL PRIMARY KEY,
    lable text NOT NULL,
    model text NOT NULL,
    year int NOT NULL,
    available bool NOT NULL DEFAULT true
);
