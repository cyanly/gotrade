CREATE TABLE IF NOT EXISTS person
(
  person_id       SERIAL,
  user_id         TEXT NOT NULL,
  organisation_id INT,
  first_name      TEXT NOT NULL,
  last_name       TEXT NOT NULL,
  title           TEXT NOT NULL,
  role_type       TEXT NOT NULL,

  PRIMARY KEY (person_id)
);

CREATE TABLE IF NOT EXISTS organisation
(
  organisation_id      SERIAL,
  name                 TEXT UNIQUE NOT NULL,
  marketing_name       TEXT        NOT NULL,
  organisation_type    TEXT        NOT NULL,

  PRIMARY KEY (organisation_id)
);

