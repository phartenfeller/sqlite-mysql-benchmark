# sqlite-mysql-benchmark

## Generate Sample Data

Requires [node.js](https://nodejs.org/)

```sh
cd data-factory
npm install
node index.js
```

## Queries

All posts of a certain tag

```sql
select * 
  from blog_posts p
  join blog_post_tags pt
    on p.post_id = pt.post_id
 where pt.tag_id = 30
```
