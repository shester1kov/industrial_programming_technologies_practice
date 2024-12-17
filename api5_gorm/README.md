# Практическая работа API4, Шестериков Дмитрий, ЭФМО-01-24
# Тема: Обработка ошибок, пагинация и фильтрация данных в REST API

## Примеры запросов

### GET ../products/price-range?minPrice=100&maxPrice=1000
![image](https://github.com/user-attachments/assets/47e67bb6-447c-4b36-83e9-fd3a4aa1b4d2)

### GET ../products/count-by-manufacturer
![image](https://github.com/user-attachments/assets/52fbbbcc-e9ac-424e-8ac5-4042c33db2dd)

### PUT ../products/manufacturer?manufacturer=NewManufacturer
![image](https://github.com/user-attachments/assets/f6122225-7647-45d1-a42f-cdf86ca667ec)
![image](https://github.com/user-attachments/assets/1088fd05-5073-466d-9075-54d44223145f)



```sql
ALTER TABLE public.products ADD COLUMN price DECIMAL(10, 2);

UPDATE public.products SET price = 1200 WHERE id = 2; -- Казеиновый протеин
UPDATE public.products SET price = 1300 WHERE id = 3; -- Растительный протеин
UPDATE public.products SET price = 1500 WHERE id = 4; -- Изолят сывороточного протеина
UPDATE public.products SET price = 700 WHERE id = 5;  -- Протеиновые батончики
UPDATE public.products SET price = 900 WHERE id = 6;  -- Креатин моногидрат
UPDATE public.products SET price = 1100 WHERE id = 7; -- Креатин HCL
UPDATE public.products SET price = 1000 WHERE id = 8; -- Креатиновые капсулы
UPDATE public.products SET price = 1150 WHERE id = 9; -- Креатин с добавками
UPDATE public.products SET price = 1200 WHERE id = 10; -- Креатиновый комплекс
UPDATE public.products SET price = 950 WHERE id = 11; -- BCAA 2:1:1
UPDATE public.products SET price = 1050 WHERE id = 12; -- BCAA с электролитами
UPDATE public.products SET price = 850 WHERE id = 13; -- Глютамин
UPDATE public.products SET price = 1100 WHERE id = 14; -- Комплекс EAA
UPDATE public.products SET price = 800 WHERE id = 15; -- Аминокислоты в таблетках
UPDATE public.products SET price = 1300 WHERE id = 16; -- Мультивитамины
UPDATE public.products SET price = 900 WHERE id = 17; -- Витамин D3
UPDATE public.products SET price = 750 WHERE id = 18; -- Омега-3
UPDATE public.products SET price = 500 WHERE id = 19; -- Магний и цинк
UPDATE public.products SET price = 600 WHERE id = 20; -- Антиоксиданты
UPDATE public.products SET price = 1100 WHERE id = 21; -- Сжигатель жира
UPDATE public.products SET price = 450 WHERE id = 22; -- Изотонический напиток
UPDATE public.products SET price = 1000 WHERE id = 23; -- L-карнитин
UPDATE public.products SET price = 1300 WHERE id = 24; -- Предтренировочный комплекс
UPDATE public.products SET price = 1400 WHERE id = 25; -- Коллаген
UPDATE public.products SET price = 1300 WHERE id = 1; -- Сывороточный протеин

alter table products
add column manufacturer varchar(255) not null default 'unknown';


UPDATE products
SET manufacturer = CASE
    WHEN id < 10 THEN 'manufacturer1'
    ELSE 'manufacturer2'
END;


```
