CREATE TABLE blog_posts (
  post_id INTEGER not null primary key AUTOINCREMENT
, content TEXT
, created timestamp default CURRENT_TIMESTAMP
);

CREATE TABLE blog_tags (
  tag_id INTEGER not null primary key AUTOINCREMENT
, descr  varchar(100)
);

CREATE TABLE blog_post_tags (
  post_id INTEGER not null
, tag_id INTEGER not null
, primary key (post_id, tag_id)
, FOREIGN KEY (tag_id) references blog_tags (tag_id)
, FOREIGN KEY (post_id) references blog_posts (post_id)
);
