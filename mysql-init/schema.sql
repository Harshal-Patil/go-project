CREATE DATABASE IF NOT EXISTS mydatabase;

USE mydatabase;

CREATE TABLE IF NOT EXISTS employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    company_name VARCHAR(100),
    address VARCHAR(255),
    city VARCHAR(100),
    county VARCHAR(100),
    postal VARCHAR(10),
    phone VARCHAR(20),
    email VARCHAR(255),
    web VARCHAR(1000)
);
