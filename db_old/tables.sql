CREATE TABLE blog_posts (
  post_id    INTEGER not null primary key AUTOINCREMENT
, content    TEXT         not null
, title      varchar(200) not null
, slug       varchar(200) not null
, created_dt datetime
, unique (slug)
);

CREATE TABLE blog_tags (
  tag_id INTEGER not null primary key AUTOINCREMENT
, descr  varchar(100) not null
);

CREATE TABLE blog_post_tags (
  post_id INTEGER not null
, tag_id  INTEGER not null
, primary key (post_id, tag_id)
, FOREIGN KEY (tag_id) references blog_tags (tag_id)
, FOREIGN KEY (post_id) references blog_posts (post_id)
);

CREATE TABLE blog_comments (
  comment_id INTEGER not null primary key AUTOINCREMENT
, post_id    INTEGER not null
, user_name  varchar(100) not null
, user_email varchar(100) not null
, comment    text not null
, FOREIGN KEY (post_id) references blog_posts (post_id)
);
