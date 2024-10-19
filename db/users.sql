CREATE TABLE users (
  id       INT AUTO_INCREMENT NOT NULL,
  username VARCHAR(255) UNIQUE NOT NULL,
  pwhash   VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO users
  (username, pwhash)
VALUES
  ('testuser', 'testpass');
