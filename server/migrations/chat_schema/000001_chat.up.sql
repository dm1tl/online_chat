CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255)
);
CREATE TABLE clients (
    id INT UNIQUE,
    username VARCHAR(255) NOT NULL,
    room_id INT,
    PRIMARY KEY (id, room_id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    client_id INT,
    room_id INT,
    content VARCHAR(255),
    FOREIGN KEY (client_id, room_id) REFERENCES clients(id, room_id)
);