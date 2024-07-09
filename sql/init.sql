DROP DATABASE IF EXISTS explore;
CREATE DATABASE explore;
CREATE USER 'exploreuser'@'localhost' IDENTIFIED BY 'test';
GRANT ALL ON explore.* TO 'exploreuser'@'localhost';
USE explore;

DROP TABLE IF EXISTS gender;
CREATE TABLE gender (
    ID int NOT NULL auto_increment,
    Name char (16) NOT NULL,
    PRIMARY KEY (ID)
);
INSERT INTO gender ( ID, Name )
VALUES
(1, "Male"),
(2, "Female");

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    ID int NOT NULL auto_increment,
    FirstName varchar (255),
    LastName varchar (255),
    GenderID int NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (GenderID) REFERENCES gender(ID)
);
INSERT INTO users ( FirstName, LastName, GenderID )
VALUES
("John", "Stone", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Ponnappa", "Priya", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Mia", "Wong", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Peter", "Stanbridge", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Natalie", "Lee-Walsh", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Ang", "Li", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Nguta", "Ithya", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Tamzyn", "French", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Salome", "Simoes", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Trevor", "Virtue", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Tarryn", "Campbell-Gillies", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Eugenia", "Anders", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Andrew", "Kazantzis", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Verona", "Blair", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Jane", "Meldrum", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Maureen", "M. Smith", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Desiree", "Burch", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1)),
("Daly", "Harry", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Hayman", "Andrews", (SELECT ID FROM gender WHERE Name="Male" LIMIT 1)),
("Ruveni", "Ellawala", (SELECT ID FROM gender WHERE Name="Female" LIMIT 1));

CREATE TABLE decisions (
    ID int NOT NULL auto_increment,
    ActorID int NOT NULL,
    RecipientID int NOT NULL,
    Liked int NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (ActorID) REFERENCES users(ID),
    FOREIGN KEY (RecipientID) REFERENCES users(ID)
);

INSERT INTO decisions ( ActorID, RecipientID, Liked )
VALUES 
(
    (SELECT ID FROM users WHERE FirstName="John" AND LastName="Stone" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
1),
(
    (SELECT ID FROM users WHERE FirstName="Ruveni" AND LastName="Ellawala" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
1),
( 
    (SELECT ID FROM users WHERE FirstName="Peter" AND LastName="Stanbridge" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
1),
( 
    (SELECT ID FROM users WHERE FirstName="Salome" AND LastName="Simoes" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
1),
( 
    (SELECT ID FROM users WHERE FirstName="Hayman" AND LastName="Andrews" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
1),
(
    (SELECT ID FROM users WHERE FirstName="Daly" AND LastName="Harry" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
1),
(
    (SELECT ID FROM users WHERE FirstName="Andrew" AND LastName="Kazantzis" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
1),
( 
    (SELECT ID FROM users WHERE FirstName="Ponnappa" AND LastName="Priya" LIMIT 1),
    (SELECT ID FROM users WHERE FirstName="John" AND LastName="Stone" LIMIT 1),
1);
