CREATE TABLE friendship (
  id SERIAL PRIMARY KEY,
  user_id INT,
  friend_id INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (friend_id) REFERENCES users(id)
);

CREATE INDEX idx_friendship_user_id ON friendship (user_id);
CREATE INDEX idx_friendship_friend_id ON friendship (friend_id);
