CREATE TABLE USERS (
    name VARCHAR(31) NOT NULL,
    surname VARCHAR(31) NOT NULL,
    id VARCHAR(31) NOT NULL PRIMARY KEY,
    age INT NOT NULL,
    sex VARCHAR(31) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

INSERT INTO USERS (name, surname, id) VALUES ('Alipasha', 'Montaseri', '99109999');
INSERT INTO USERS (name, surname, id) VALUES ('Yasmin', 'Kadkhodaee', '98109898');
