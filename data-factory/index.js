const faker = require("faker");
const fs = require("fs");

const POST_COUNT = 500;
const TAG_COUNT = 50;
const POST_TAGS = POST_COUNT * 3;
const COMMENTS = POST_COUNT * 10;

const randomNumberOneTo = (max) => Math.ceil(Math.random() * max);

const formatDate = (date) => {
  return date
    .toISOString()
    .replace(/T/, " ") // replace T with a space
    .replace(/\..+/, ""); // delete the dot and everything after
};

const posts = [];
const tags = [];
const postTags = [];
const comments = [];

const postTagCombos = [];

// blogposts
for (let i = 1; i <= POST_COUNT; i++) {
  const titleWords = randomNumberOneTo(15) + 10;
  const title = faker.lorem.words(titleWords).replace(/'/g, "`");
  const slug = faker.unique(faker.lorem.slug).replace(/'/g, "`");
  const created = formatDate(faker.date.past(10));

  const paragraphs = randomNumberOneTo(10);
  let text = "";
  for (let j = 1; j <= paragraphs; j++) {
    text += faker.lorem.paragraphs().replace(/'/g, "`");
  }
  const insert = `insert into blog_posts (post_id, content, title, slug, created_dt) values (${i}, '${text}',\n '${title}', '${slug}', '${created}');`;
  posts.push(insert);
}

// tags
for (let i = 1; i <= TAG_COUNT; i++) {
  const tag = faker.unique(faker.animal.bird).replace(/'/g, "`");

  const insert = `insert into blog_tags (tag_id, descr) values (${i}, '${tag}');`;
  tags.push(insert);
}

// post tags
for (let i = 1; i <= POST_TAGS; i++) {
  let tagId;
  let postId;
  let unique = false;
  while (!unique) {
    tagId = randomNumberOneTo(TAG_COUNT);
    postId = randomNumberOneTo(POST_COUNT);
    const id = `${tagId}-${postId}`;
    if (!postTagCombos.includes(id)) {
      postTagCombos.push(id);
      unique = true;
    }
  }

  const insert = `insert into blog_post_tags (post_id, tag_id) values (${postId}, ${tagId});`;
  postTags.push(insert);
}

// comments
for (let i = 1; i <= COMMENTS; i++) {
  const postId = randomNumberOneTo(POST_COUNT);

  const username = faker.internet.userName().replace(/'/g, "`");
  const email = faker.internet.email().replace(/'/g, "`");

  const sentences = randomNumberOneTo(10);
  let comment = "";

  for (let j = 1; j <= sentences; j++) {
    comment += faker.lorem.sentence().replace(/'/g, "`");
  }

  const insert = `insert into blog_comments (post_id, user_name, user_email, comment) values (${postId}, '${username}', '${email}', '${comment}');`;
  comments.push(insert);
}

let fullScript = "/* Blog Posts */\n";
fullScript += posts.join("\n\n\n");
fullScript += "\n\n\n\n\n";
fullScript += "/* Tags */\n";
fullScript += tags.join("\n\n\n");
fullScript += "\n\n\n\n\n";
fullScript += "/* Post Tags */\n";
fullScript += postTags.join("\n\n\n");
fullScript += "\n\n\n\n\n";
fullScript += "/* Comments */\n";
fullScript += comments.join("\n\n\n");

fs.writeFileSync("../db/inserts.sql", fullScript);
