-- Восстанавливаем колонку rating
ALTER TABLE movies ADD COLUMN rating FLOAT;
ALTER TABLE series ADD COLUMN rating FLOAT;