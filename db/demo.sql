-- Создание схемы
-- Марки автомобилей
CREATE TABLE IF NOT EXISTS brands (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

-- Модели автомобилей
CREATE TABLE IF NOT EXISTS models (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

-- Автомобили
CREATE TABLE IF NOT EXISTS cars (
  id SERIAL PRIMARY KEY,
  brand_id INT NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
  model_id INT NOT NULL REFERENCES models(id) ON DELETE CASCADE,
  year INT NOT NULL,
  price NUMERIC(10,1) NOT NULL     -- цена в тысячах рублей
);

-- Сброс (на время разработки)
TRUNCATE TABLE cars RESTART IDENTITY CASCADE;
TRUNCATE TABLE brands RESTART IDENTITY CASCADE;
TRUNCATE TABLE models RESTART IDENTITY CASCADE;

-- Вставка 12 брендов
INSERT INTO brands (name) VALUES
('Toyota'), ('Hyundai'), ('Volkswagen'), ('BMW'),
('Mercedes-Benz'), ('Kia'), ('Audi'), ('Lada'),
('Nissan'), ('Renault'), ('Mazda'), ('Ford');

-- Вставка 12 моделей (упрощённые примеры)
INSERT INTO models (name) VALUES
('Camry'), ('Solaris'), ('Polo'), ('3 Series'),
('C-Class'), ('Rio'), ('A4'), ('Vesta'),
('Qashqai'), ('Logan'), ('CX-5'), ('Focus');

-- Вставка 12 автомобилей (со связями 1..12)
INSERT INTO cars (brand_id, model_id, year, price) VALUES
(1, 1, 2018, 1500.0),
(2, 2, 2020, 950.0),
(3, 3, 2019, 900.0),
(4, 4, 2017, 1800.0),
(5, 5, 2021, 2500.0),
(6, 6, 2020, 920.0),
(7, 7, 2018, 1600.0),
(8, 8, 2022, 850.0),
(9, 9, 2021, 1400.0),
(10, 10, 2017, 700.0),
(11, 11, 2019, 1700.0),
(12, 12, 2016, 600.0);

-- Индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_brands_name ON brands (name);
CREATE INDEX IF NOT EXISTS idx_models_name ON models (name);
CREATE INDEX IF NOT EXISTS idx_cars_year ON cars (year);
CREATE INDEX IF NOT EXISTS idx_cars_price ON cars (price);
