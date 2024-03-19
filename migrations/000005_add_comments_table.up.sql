CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  post_id INT,
  commentator_id INT,
  comment TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (commentator_id) REFERENCES users(id)
);

CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_commentator_id ON comments (commentator_id);
