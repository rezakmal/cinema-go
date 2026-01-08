-- Create database (run this manually)
-- CREATE DATABASE cinema_db;

-- Create cinema table
CREATE TABLE IF NOT EXISTS cinema_db (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location TEXT NOT NULL,
    rating DECIMAL(3,1) DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

-- Create index for better performance
CREATE INDEX idx_cinema_name ON cinema(name);
CREATE INDEX idx_cinema_location ON cinema(location);