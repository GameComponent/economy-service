CREATE TABLE IF NOT EXISTS item (  
  id UUID DEFAULT gen_random_uuid() NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  name STRING DEFAULT '' NOT NULL,
  metadata JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id)
);

CREATE INDEX index_name ON item(name);

CREATE TABLE IF NOT EXISTS currency (  
  id UUID DEFAULT gen_random_uuid() NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  name STRING DEFAULT '' NOT NULL UNIQUE,
  short_name STRING DEFAULT '' NOT NULL UNIQUE,
  symbol STRING DEFAULT '' NOT NULL UNIQUE,
  metadata JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS player (
	id STRING NOT NULL,
	name STRING NOT NULL,

	PRIMARY KEY (id)
);

CREATE INDEX index_name ON player(name);

CREATE TABLE IF NOT EXISTS storage (
  id UUID DEFAULT gen_random_uuid() NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  player_id STRING NOT NULL,
  name STRING NOT NULL,
  metadata JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id),
	FOREIGN KEY (player_id) REFERENCES player(id)
);

CREATE INDEX index_player_id ON storage(player_id);

CREATE TABLE IF NOT EXISTS storage_item (
  id UUID DEFAULT gen_random_uuid() NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  item_id UUID NOT NULL,
  storage_id UUID NOT NULL,
  metadata JSONB DEFAULT '{}' NOT null,
  
  PRIMARY KEY (id),
  FOREIGN KEY (item_id) REFERENCES item(id),
  FOREIGN KEY (storage_id) REFERENCES storage(id)
);

CREATE TABLE IF NOT EXISTS storage_currency (
  id UUID DEFAULT gen_random_uuid() NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  currency_id UUID NOT NULL,
  storage_id UUID NOT NULL,
  amount INT64 DEFAULT 0 NOT NULL,
  
  PRIMARY KEY (id),
  FOREIGN KEY (currency_id) REFERENCES currency(id),
  FOREIGN KEY (storage_id) REFERENCES storage(id),
  UNIQUE (currency_id, storage_id)
);

CREATE TABLE IF NOT EXISTS config (  
  key STRING NOT NULL,
	value JSONB,
  
  PRIMARY KEY (key)
);

CREATE TABLE IF NOT EXISTS account (  
  id UUID DEFAULT gen_random_uuid() NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp(),
	email STRING NOT NULL,
	password STRING NOT NULL,
  
  PRIMARY KEY (id)
);
