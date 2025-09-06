drop table if exists users, posts, comments;

create table users (
id SERIAL primary key,
username VARCHAR(32) unique not null,
email varchar(128) unique not null,
created_at TIMESTAMP default NOW()
);

create table posts (
id SERIAL primary key,
title varchar(64) not null,
content TEXT not null, 
published_at timestamp default now(),

author_id int,

CONSTRAINT fk_posts_users
    foreign key (author_id)
    REFERENCES users(id)
);

create table comments (
id SERIAL primary key,
comment_text TEXT not null,
created_at TIMESTAMP DEFAULT NOW(),

post_id INT,
user_id INT,

CONSTRAINT fk_comments_posts
    foreign key (post_id)
    REFERENCES posts(id),

CONSTRAINT fk_comments_users
    foreign key (user_id)
    REFERENCES users(id)
);

--
-- CREATE PROCEDURES!!!
--

create procedure CreateUser (u varchar, e varchar)
LANGUAGE sql AS $$
INSERT INTO users(username,email) VALUES (u,e);
$$;

-- example
-- CALL CreateUser('golang_lover', 'haskell@husking.hosk');