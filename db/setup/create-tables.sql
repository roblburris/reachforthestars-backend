-- Creates necessary Tables, only meant to be run initially
CREATE TABLE BLOG_POSTS (
    blogid INT PRIMARY KEY,
    author VARCHAR(100),
    datePosted VARCHAR(10),
    duration INT,
    url BYTEA,
    content BYTEA
);

-- CREATE TABLE USERS (
--     uid INT PRIMARY KEY,
--     name VARCHAR(100),
--     bio BYTEA
-- );