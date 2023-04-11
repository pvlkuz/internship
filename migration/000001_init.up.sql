
CREATE EXTENSION "uuid-ossp";

CREATE TABLE IF NOT EXISTS records
(
     id UUID PRIMARY KEY, 
     transform_type TEXT NOT NULL,
     caesar_shift INT,
     result TEXT,
     created_at TIMESTAMP DEFAULT now(),
     updated_at TIMESTAMP DEFAULT now()
);

CREATE OR REPLACE FUNCTION update_stamp() RETURNS trigger AS $update_stamp$
BEGIN 
     NEW.updated_at = now();
     RETURN NEW;
END;

$update_stamp$ language plpgsql;

CREATE TRIGGER update_time BEFORE UPDATE ON records FOR EACH ROW EXECUTE PROCEDURE update_stamp();
