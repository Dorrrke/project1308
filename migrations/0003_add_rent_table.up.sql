CREATE TABLE IF NOT EXISTS rents (
    rid varchar(36) NOT NULL PRIMARY KEY,
    car_id varchar(36) NOT NULL,
    user_id varchar(36) NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    FOREIGN KEY (car_id) REFERENCES cars(cid),
    FOREIGN KEY (user_id) REFERENCES users(uid)
)