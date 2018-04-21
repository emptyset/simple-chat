CREATE TABLE user (
	id INT AUTO_INCREMENT NOT NULL,
	username VARCHAR(255),
	hash BINARY(60),
	PRIMARY KEY (id)
);

CREATE TABLE message (
	id INT AUTO_INCREMENT NOT NULL,
	timestamp INT NOT NULL,
	sender_id INT NOT NULL,
	recipient_id INT NOT NULL,
	content TEXT,
	metadata TEXT,

	PRIMARY KEY (id),
	FOREIGN KEY sender_id REFERENCES user(id),
	FOREIGN KEY recipient_id REFERENCES user(id)
);
