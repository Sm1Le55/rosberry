DROP SCHEMA IF EXISTS rosberry_fsm CASCADE;
DROP USER IF EXISTS rosberry_fsm;
CREATE USER rosberry_fsm WITH ENCRYPTED PASSWORD '123';
CREATE SCHEMA AUTHORIZATION rosberry_fsm;

SET schema 'rosberry_fsm';

CREATE TABLE Users (
    ID serial NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
	password varchar(128) NOT NULL,
    accessKey varchar(128),
    accessKeyExpireDate timestamp,
    CONSTRAINT User_pk PRIMARY KEY (ID)
) WITH (
  OIDS=FALSE
);

CREATE TABLE Profile (
	ID serial NOT NULL,
    userID integer NOT NULL,
	name varchar(255) NOT NULL,
	photo bytea,
	birthday DATE NOT NULL,
	CONSTRAINT Profile_pk PRIMARY KEY (ID)
) WITH (
  OIDS=FALSE
);

CREATE TABLE Locations (
    ID serial NOT NULL,
    title varchar(32) NOT NULL UNIQUE,
    CONSTRAINT Location_pk PRIMARY KEY (ID)
) WITH (
  OIDS=FALSE
);

CREATE TABLE LocationSettings (
    userID integer NOT NULL UNIQUE,
    location integer NOT NULL
) WITH (
  OIDS=FALSE
);

CREATE TABLE AuthHistory (
	ID serial NOT NULL,
	userID integer NOT NULL,
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	latitude FLOAT NOT NULL,
	longitude FLOAT NOT NULL,
	CONSTRAINT AuthHistory_pk PRIMARY KEY (ID)
) WITH (
  OIDS=FALSE
);

CREATE TABLE AgeSettings (
	userID integer NOT NULL UNIQUE,
	showRangeForMe integer NOT NULL,
	hideMeByRange integer NOT NULL
) WITH (
  OIDS=FALSE
);

CREATE TABLE AgeRanges (
	ID serial NOT NULL,
	title varchar(128) NOT NULL,
	minAge integer NOT NULL,
	maxAge integer,
	CONSTRAINT AgeRanges_pk PRIMARY KEY (ID)
) WITH (
  OIDS=FALSE
);

CREATE TABLE ShowInterestsSettings (
	userID integer NOT NULL,
	theme integer NOT NULL
) WITH (
  OIDS=FALSE
);

CREATE TABLE Themes (
	ID serial NOT NULL,
	title varchar(255) NOT NULL UNIQUE,
	CONSTRAINT Themes_pk PRIMARY KEY (ID)
) WITH (
  OIDS=FALSE
);

CREATE TABLE HideInterestsSettings (
	userID integer NOT NULL,
	theme integer NOT NULL
) WITH (
  OIDS=FALSE
);

CREATE TABLE UserInterest (
	userID integer NOT NULL,
	theme integer NOT NULL
) WITH (
  OIDS=FALSE
);

ALTER TABLE Profile ADD CONSTRAINT Profile_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);
ALTER TABLE AuthHistory ADD CONSTRAINT AuthHistory_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);

ALTER TABLE LocationSettings ADD CONSTRAINT AgeSettings_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);
ALTER TABLE LocationSettings ADD CONSTRAINT AgeSettings_fk1 FOREIGN KEY (location) REFERENCES Locations(ID);

ALTER TABLE AgeSettings ADD CONSTRAINT AgeSettings_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);
ALTER TABLE AgeSettings ADD CONSTRAINT AgeSettings_fk1 FOREIGN KEY (showRangeForMe) REFERENCES AgeRanges(ID);
ALTER TABLE AgeSettings ADD CONSTRAINT AgeSettings_fk2 FOREIGN KEY (hideMeByRange) REFERENCES AgeRanges(ID);

ALTER TABLE ShowInterestsSettings ADD CONSTRAINT ShowInterestsSettings_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);
ALTER TABLE ShowInterestsSettings ADD CONSTRAINT ShowInterestsSettings_fk1 FOREIGN KEY (theme) REFERENCES Themes(ID);

ALTER TABLE HideInterestsSettings ADD CONSTRAINT HideInterestsSettings_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);
ALTER TABLE HideInterestsSettings ADD CONSTRAINT HideInterestsSettings_fk1 FOREIGN KEY (theme) REFERENCES Themes(ID);

ALTER TABLE ProfileInterest ADD CONSTRAINT ProfileInterest_fk0 FOREIGN KEY (userID) Users(ID);
ALTER TABLE ProfileInterest ADD CONSTRAINT ProfileInterest_fk1 FOREIGN KEY (theme) REFERENCES Themes(ID);

INSERT INTO Themes (title)
VALUES
    ('Music'),
    ('Movies'),
    ('Art'),
    ('Animals'),
    ('Little cute kittens');

INSERT INTO AgeRanges (title,minAge,maxAge)
VALUES
    ('all',0, null),
    ('18-24', 18, 24),
    ('25-40', 25, 40),
    ('40+', 41, null);

INSERT INTO Locations (title)
VALUES
    ('world'),
    ('сountry'),
    ('nearby');

INSERT INTO Users (email,password,accessKey,accessKeyExpireDate)
VALUES
    ('User1@mail.ru','123','55555','2020-01-30 09:10:13.65472'),
    ('User2@mail.ru','234','55555','2020-05-30 09:10:13.65472'),
    ('User3@mail.ru','345','55555','2020-05-30 09:10:13.65472'),
    ('User4@mail.ru','456','55555','2020-05-30 09:10:13.65472');
    
INSERT INTO Profile (userID,name,birthday)
VALUES
    (1, 'Mike','21.12.1979'),
    (2, 'Anna','01.01.2000'),
    (3, 'Dave','10.06.1977'),
    (4, 'Katy','16.08.1994');

INSERT INTO LocationSettings (profile, location)
VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (4, 1);

INSERT INTO AgeSettings (profile, showRangeForMe, hideMeByRange)
VALUES
    (1, 2, 3),
    (2, 3, 4),
    (3, 2, 1),
    (4, 1, 1);

INSERT INTO ProfileInterest (profile, theme)
VALUES
    (1, 1),
    (1, 3),
    (2, 1),
    (2, 4),
    (2, 5),
    (3, 3),
    (3, 2),
    (4, 4);

INSERT INTO ShowInterestsSettings (profile, theme)
VALUES
    (4, 1),
    (4, 3),
    (3, 1),
    (3, 4),
    (3, 5),
    (2, 3),
    (1, 2),
    (1, 4);

INSERT INTO HideInterestsSettings (profile, theme)
VALUES
    (4, 4),
    (4, 2),
    (3, 2),
    (2, 3),
    (2, 1),
    (2, 3),
    (1, 3),
    (1, 1);