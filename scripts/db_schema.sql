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
-- using functions, since procedures are a pain to return something
-- prepare queries only live for a certain time in the current session

drop function if exists CreateUser;

CREATE FUNCTION CreateUser (_username VARCHAR, _email VARCHAR)
 RETURNS TABLE (
  id INTEGER,
  username VARCHAR,
  email VARCHAR
 ) -- we specify the type of what we return. could be SETOF users if we return the entire table
  LANGUAGE plpgsql  AS
$func$
BEGIN
	RETURN QUERY -- we return this result
      INSERT INTO users(username, email)
      VALUES (_username, _email)
      RETURNING username,email;
END
$func$;

SELECT * FROM CreateUser('angular_supporter12', 'htmx12@htmx.com');

