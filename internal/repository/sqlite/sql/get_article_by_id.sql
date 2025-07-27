SELECT 
  id, 
  title,
  author, 
  body, 
  published_at, 
  updated_at 
FROM ARTICLES 
WHERE id = ?
