-- Удаляем колонку rating из таблицы movies
ALTER TABLE movies DROP COLUMN IF EXISTS rating;

-- Удаляем колонку rating из таблицы series
ALTER TABLE series DROP COLUMN IF EXISTS rating;