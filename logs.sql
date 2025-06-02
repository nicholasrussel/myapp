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
