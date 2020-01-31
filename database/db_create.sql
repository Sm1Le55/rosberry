DROP SCHEMA IF EXISTS rosberry_fsm CASCADE;
DROP USER rosberry_fsm;
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
	profile integer NOT NULL,
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
	profile integer NOT NULL,
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
	profile integer NOT NULL,
	theme integer NOT NULL
) WITH (
  OIDS=FALSE
);

CREATE TABLE ProfileInterest (
	profile integer NOT NULL,
	theme integer NOT NULL
) WITH (
  OIDS=FALSE
);

ALTER TABLE Profile ADD CONSTRAINT Profile_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);
ALTER TABLE AuthHistory ADD CONSTRAINT AuthHistory_fk0 FOREIGN KEY (userID) REFERENCES Users(ID);

ALTER TABLE AgeSettings ADD CONSTRAINT AgeSettings_fk0 FOREIGN KEY (profile) REFERENCES Profile(ID);
ALTER TABLE AgeSettings ADD CONSTRAINT AgeSettings_fk1 FOREIGN KEY (showRangeForMe) REFERENCES AgeRanges(ID);
ALTER TABLE AgeSettings ADD CONSTRAINT AgeSettings_fk2 FOREIGN KEY (hideMeByRange) REFERENCES AgeRanges(ID);

ALTER TABLE ShowInterestsSettings ADD CONSTRAINT ShowInterestsSettings_fk0 FOREIGN KEY (profile) REFERENCES Profile(ID);
ALTER TABLE ShowInterestsSettings ADD CONSTRAINT ShowInterestsSettings_fk1 FOREIGN KEY (theme) REFERENCES Themes(ID);

ALTER TABLE HideInterestsSettings ADD CONSTRAINT HideInterestsSettings_fk0 FOREIGN KEY (profile) REFERENCES Profile(ID);
ALTER TABLE HideInterestsSettings ADD CONSTRAINT HideInterestsSettings_fk1 FOREIGN KEY (theme) REFERENCES Themes(ID);

ALTER TABLE ProfileInterest ADD CONSTRAINT ProfileInterest_fk0 FOREIGN KEY (profile) REFERENCES Profile(ID);
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
    ('beetwen 18 and 24', 18, 24),
    ('beetwen 25 and 40', 25, 40),
    ('over 40', 41, null);
    
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