SELECT
  a.id,
  a.title,
  a.body,
  u.user_name AS "author.user_name",
  u.first_name AS "author.first_name",
  u.last_name AS "author.last_name"
FROM articles a
JOIN users u ON a.author = u.user_name
ORDER BY a.id ASC
