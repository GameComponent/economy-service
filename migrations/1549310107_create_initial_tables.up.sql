CREATE TABLE IF NOT EXISTS storage (
  id UUID DEFAULT gen_random_uuid() NOT NULL,
  player_id STRING NOT NULL UNIQUE,
  name STRING NOT NULL,
  data JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id)
);

CREATE INDEX index_player_id ON storage(player_id);

CREATE TABLE IF NOT EXISTS item (  
  id UUID DEFAULT gen_random_uuid() NOT NULL,
  name STRING DEFAULT '' NOT NULL,
  data JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS storage_item (
  id UUID DEFAULT gen_random_uuid() NOT NULL,
  item_id UUID NOT NULL,
  storage_id UUID NOT NULL,
  data JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id),
  FOREIGN KEY (item_id) REFERENCES item (id),
  FOREIGN KEY (storage_id) REFERENCES storage (id)
);
