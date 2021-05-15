const faker = require("faker");
const fs = require("fs");

const posts = [];

// blogposts
for (let i = 1; i <= 5000; i++) {
  const paragraphs = Math.ceil(Math.random() * 10);

  let text = "";

  for (let j = 1; j <= paragraphs; j++) {
    text += faker.lorem.paragraphs().replace(/'/g, "`");
  }
  const insert = `insert into blog_posts (post_id, content) values (${i}, '${text}');`;
  posts.push(insert);
}

fs.writeFileSync("../db/inserts.sql", posts.join("\n\n\n"));
