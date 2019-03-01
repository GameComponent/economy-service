CREATE TABLE IF NOT EXISTS notification (
  id UUID DEFAULT gen_random_uuid() NOT NULL,
 	player_id STRING NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  metadata JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id),
  FOREIGN KEY (player_id) REFERENCES player (id)
);
