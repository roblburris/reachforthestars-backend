-- Creates necessary Tables, only meant to be run initially
CREATE TABLE BLOG_POSTS (
    blogid INT,
    author VARCHAR(100),
    datePosted VARCHAR(10),
    duration INT,
    url BYTEA,
    content BYTEA,
    PRIMARY KEY(blogid)
);

CREATE TABLE BLOG_POST_TITLES (
    blogid INT,
    blogTitle VARCHAR(100),
    FOREIGN KEY(blogID) REFERENCES BLOG_POSTS(blogid)
);
-- CREATE TABLE USERS (
--     uid INT PRIMARY KEY,
--     name VARCHAR(100),
--     bio BYTEA
-- );