-- Active: 1779046001743@@localhost@5432@gratedata

CREATE TABLE IF NOT EXISTS users (
    idme UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    id VARCHAR(255) NOT NULL,
    email VARCHAR(255)NOT NULL,
    verified_email BOOLEAN ,
    name VARCHAR(255),
    picture VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);