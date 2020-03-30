
CREATE TABLE IF NOT EXISTS pages (
   id serial PRIMARY KEY,
   url VARCHAR (150) NOT NULL,
   content TEXT,
   created_date date,
   created_at TIMESTAMP
);
