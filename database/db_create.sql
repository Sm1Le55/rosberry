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
    userID integer NOT NULL,
	name varchar(255) NOT NULL,
	photo bytea,
	birthday DATE NOT NULL,
    country varchar,
	CONSTRAINT Profile_pk PRIMARY KEY (userID)
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
	coord point NOT NULL,
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

CREATE TABLE Empty (
    Avatar bytea
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

ALTER TABLE UserInterest ADD CONSTRAINT UserInterest_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);
ALTER TABLE UserInterest ADD CONSTRAINT UserInterest_fk1 FOREIGN KEY (theme) REFERENCES Themes(ID);

INSERT INTO Empty (Avatar)
VALUES 
(
    decode(
    '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'
    ,'base64')
);

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
    
INSERT INTO Profile (userID,name,birthday,country)
VALUES
    (1, 'Mike','21.12.1979', 'USA'),
    (2, 'Anna','01.01.2000', 'Russia'),
    (3, 'Dave','10.06.1977', 'Russia'),
    (4, 'Katy','16.08.1994', 'China');

INSERT INTO LocationSettings (userID, location)
VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (4, 1);

INSERT INTO AgeSettings (userID, showRangeForMe, hideMeByRange)
VALUES
    (1, 2, 3),
    (2, 3, 4),
    (3, 2, 1),
    (4, 1, 1);

INSERT INTO UserInterest (userID, theme)
VALUES
    (1, 1),
    (1, 3),
    (2, 1),
    (2, 4),
    (2, 5),
    (3, 3),
    (3, 2),
    (4, 4);

INSERT INTO ShowInterestsSettings (userID, theme)
VALUES
    (4, 1),
    (4, 3),
    (3, 1),
    (3, 4),
    (3, 5),
    (2, 3),
    (1, 2),
    (1, 4);

INSERT INTO HideInterestsSettings (userID, theme)
VALUES
    (4, 1),
    (2, 3),
    (2, 1),
    (2, 3),
    (1, 3),
    (1, 1);

INSERT INTO authhistory(userid, "time", coord)
VALUES 
(1, '2020-01-01 04:44:17.519673', '55.7455, 59.1523'),
(1, '2020-02-01 08:04:17.519673', '49.1215, 88.1227'),
(2, '2020-02-01 10:04:17.519673', '55.5615, 59.5527'),
(3, '2020-02-01 08:37:22.519673', '33.2215, 72.1224'),
(4, '2020-02-02 04:04:17.519673', '47.3155, 24.1587');


CREATE FUNCTION fdistance(src_lat double precision, src_lon double precision, dst_lat double precision, dst_lon double precision) RETURNS double precision
    LANGUAGE plpgsql
    AS $$
DECLARE
Distance double precision;
BEGIN

Distance = 6371 * 2 * ASIN(SQRT(
    POWER(SIN((src_lat - ABS(dst_lat)) * PI()/180 / 2), 2) + 
    COS(src_lat * PI()/180) * 
    COS(ABS(dst_lat) * PI()/180) * 
    POWER(SIN((src_lon - dst_lon) * PI()/180 / 2), 2)
));
RETURN Distance;
END;
$$;


CREATE FUNCTION pointDistance(p1 point, p2 point) RETURNS double precision
    LANGUAGE plpgsql
    AS $$
DECLARE
Distance double precision;
BEGIN

Distance = rosberry_fsm.fdistance(p1[0], p1[1], p2[0], p2[1]);

RETURN Distance;
END;
$$;