CREATE TABLE tags (
  id SERIAL PRIMARY KEY,
  post_id INT,
  name VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE INDEX idx_tags_post_id ON tags (post_id);
