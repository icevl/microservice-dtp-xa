CREATE TABLE IF NOT EXISTS users (
  uuid VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  PRIMARY KEY(email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8