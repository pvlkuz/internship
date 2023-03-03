
CREATE EXTENSION "uuid-ossp";

CREATE TABLE IF NOT EXISTS Records
(
     id UUID PRIMARY KEY, 
     transform_type TEXT NOT NULL,
     caesar_shift INT,
     result TEXT,
     created_at TIMESTAMP NOT NULL,
     updated_at TIMESTAMP
);
