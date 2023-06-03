CREATE TABLE USERS (
    name VARCHAR(31) NOT NULL,
    surname VARCHAR(31) NOT NULL,
    id VARCHAR(31) NOT NULL PRIMARY KEY,
    age INT,
    sex VARCHAR(31)
);

INSERT INTO USERS (name, surname, id, age, sex) 
    VALUES ('Alipasha', 'Montaseri', '99109999', 20, 'Male');
INSERT INTO USERS (name, surname, id, age, sex) 
    VALUES ('Mahdi', 'Kadkhodaee', '98109898', 21, 'Female');
