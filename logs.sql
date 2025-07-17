CREATE TABLE myapp.logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NULL, -- bisa null jika anonymous
    action_type VARCHAR(100), -- 'LOGIN', 'REGISTER', 'UPDATE_PROFILE', dll
    status VARCHAR(20), -- 'SUCCESS' atau 'FAILED'
    message TEXT,
    data TEXT, -- bisa simpan JSON ringkas request body atau response error
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO myapp.logs (user_id,action_type,status,message,`data`,ip_address,user_agent,created_at) VALUES
	 (7,'REGISTER','SUCCESS','Registrasi dan login berhasil',NULL,'192.168.65.1','PostmanRuntime/7.43.3','2025-06-02 06:53:10'),
	 (7,'TOKEN','SUCCESS','Token dan refresh token berhasil dibuat',NULL,'192.168.65.1','PostmanRuntime/7.43.3','2025-06-02 06:53:10'),
	 (3,'LOGIN','SUCCESS','Login berhasil',NULL,'192.168.65.1','PostmanRuntime/7.43.3','2025-06-02 06:53:42'),
	 (3,'TOKEN','SUCCESS','Token dan refresh token berhasil dibuat',NULL,'192.168.65.1','PostmanRuntime/7.43.3','2025-06-02 06:53:42'),
	 (NULL,'LOGIN','FAILED','Email not found',NULL,'192.168.65.1','PostmanRuntime/7.43.3','2025-06-02 06:54:03'),
	 (3,'LOGIN','SUCCESS','Login berhasil',NULL,'192.168.65.1','PostmanRuntime/7.43.3','2025-06-02 07:07:55'),
	 (3,'TOKEN','SUCCESS','Token dan refresh token berhasil dibuat',NULL,'192.168.65.1','PostmanRuntime/7.43.3','2025-06-02 07:07:55');

CREATE TABLE myapp.logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NULL, -- bisa null jika anonymous
    action_type VARCHAR(100), -- 'LOGIN', 'REGISTER', 'UPDATE_PROFILE', dll
    status VARCHAR(20), -- 'SUCCESS' atau 'FAILED'
    message TEXT,
    data TEXT, -- bisa simpan JSON ringkas request body atau response error
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE myapp.users (
  `id` int(11) NOT NULL,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `user_type` int(11) NOT NULL DEFAULT 0
)

ALTER TABLE myapp.users DROP PRIMARY KEY;

ALTER TABLE myapp.users
MODIFY COLUMN `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY;



INSERT INTO  myapp.users (`id`, `username`, `email`, `password`, `created_at`, `user_type`) VALUES
(1, 'russel', 'russel@gmail.com', '$2b$12$DdYfc.cZXX5mwoisXsm3Ae1/iDd.SLNduxlzOrYIYf3ZqVtXhlVnm\n', '2025-05-26 04:41:17', 1);

INSERT INTO  myapp.users (`id`, `username`, `email`, `password`, `created_at`, `user_type`) VALUES
(2, 'nicholas', 'nicholas@gmail.com', '$2b$12$DdYfc.cZXX5mwoisXsm3Ae1/iDd.SLNduxlzOrYIYf3ZqVtXhlVnm\n', '2025-05-26 04:41:17', 1);


CREATE TABLE myapp.messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    content TEXT NOT NULL,
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES myapp.users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES myapp.users(id) ON DELETE CASCADE
);

SELECT id, sender_id, receiver_id, content, sent_at 
FROM myapp.messages 
WHERE (sender_id = 1 AND receiver_id = 2) OR (sender_id = 2 AND receiver_id = 1)
ORDER BY sent_at ASC;

DESCRIBE myapp.messages;

ALTER TABLE myapp.messages
ADD COLUMN group_id INT NULL,
ADD CONSTRAINT fk_group_id FOREIGN KEY (group_id) REFERENCES `groups`(id) ON DELETE SET NULL;


CREATE TABLE myapp.message_receivers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    message_id INT NOT NULL,
    receiver_id INT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE myapp.groups (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  created_by INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE myapp.group_members (
  id INT PRIMARY KEY AUTO_INCREMENT,
  group_id INT NOT NULL,
  user_id INT NOT NULL,
  joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_id) REFERENCES `groups`(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


INSERT INTO myapp.message_receivers (message_id, receiver_id) VALUES (12, 3);


CREATE TABLE myapp.group_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    group_id INT NOT NULL,
    content TEXT NOT NULL,
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES myapp.users(id),
    FOREIGN KEY (group_id) REFERENCES myapp.groups(id)
);



CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user1_id INT NOT NULL,
    user2_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user1_id, user2_id)
);

CREATE TABLE friend_requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_user_id INT NOT NULL,
    to_user_id INT NOT NULL,
    status ENUM('pending', 'accepted', 'rejected') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_friend_request (from_user_id, to_user_id)
);

CREATE TABLE blocked_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    blocker_id INT NOT NULL,
    blocked_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY (blocker_id, blocked_id)
);

ALTER TABLE friends ADD COLUMN is_blocked_by INT DEFAULT NULL;

ALTER TABLE friend_requests ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

CREATE TABLE chat_rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    is_group BOOLEAN DEFAULT FALSE,
    name VARCHAR(255), -- untuk grup chat, nullable untuk 1-on-1
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chat_room_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    chat_room_id INT,
    user_id INT,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chat_room_id) REFERENCES chat_rooms(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

ALTER TABLE chat_rooms
ADD COLUMN created_by INT NULL,
ADD CONSTRAINT fk_chat_rooms_created_by FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE messages
ADD COLUMN chat_room_id INT AFTER id,
ADD CONSTRAINT fk_messages_chat_room_id FOREIGN KEY (chat_room_id) REFERENCES chat_rooms(id);




