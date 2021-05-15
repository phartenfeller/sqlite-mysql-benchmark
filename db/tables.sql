CREATE TABLE blog_posts (
  id INTEGER not null primary key AUTOINCREMENT
, content TEXT
, created timestamp default systimestamp
);
