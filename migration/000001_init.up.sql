
CREATE EXTENSION "uuid-ossp";

CREATE TABLE IF NOT EXISTS Records
(
     id UUID PRIMARY KEY, 
     transform_type TEXT NOT NULL,
     caesar_shift INT,
     result TEXT,
     created_at INT NOT NULL,
     updated_at INT
);
