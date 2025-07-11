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


