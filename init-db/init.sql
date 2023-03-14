-- create database article
CREATE DATABASE IF NOT EXISTS article;

-- use created db
USE article;

-- create table article
CREATE TABLE IF NOT EXISTS article(
    id INT PRIMARY KEY AUTO_INCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);