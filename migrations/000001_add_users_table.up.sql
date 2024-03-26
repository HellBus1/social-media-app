CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  password VARCHAR(100),
  email VARCHAR(50),
  phone VARCHAR(20) UNIQUE,
  image_url VARCHAR(100),
  credential_type VARCHAR(15),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_password ON users (password);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_phone ON users (phone);
