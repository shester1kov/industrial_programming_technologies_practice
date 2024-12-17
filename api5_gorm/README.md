# Практическая работа API5, Шестериков Дмитрий, ЭФМО-01-24
# Тема: Работа с транзакциями и сложными запросами в GORM

## Примеры запросов

### GET ../products/price-range?minPrice=100&maxPrice=1000
![image](https://github.com/user-attachments/assets/47e67bb6-447c-4b36-83e9-fd3a4aa1b4d2)

### GET ../products/count-by-manufacturer
![image](https://github.com/user-attachments/assets/52fbbbcc-e9ac-424e-8ac5-4042c33db2dd)

### PUT ../products/manufacturer?manufacturer=NewManufacturer
![image](https://github.com/user-attachments/assets/f6122225-7647-45d1-a42f-cdf86ca667ec)
![image](https://github.com/user-attachments/assets/1088fd05-5073-466d-9075-54d44223145f)

```go
package main


import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
    "strconv"
    "context"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string 
	Password string
	Role	 string

}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
	Role	 string `json:"role"`
}

func generateToken(username string, role string) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: username,
        Role:     role, // Включаем роль в токен
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}


func login(c *gin.Context) {
    var creds Credentials
    if err := c.BindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
        return
    }

    // Проверяем имя пользователя и пароль
    storedPassword, ok := users[creds.Username]
    if !ok || storedPassword != creds.Password {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
        return
    }

    // Извлекаем роль пользователя из мапы roles
    role, roleExists := roles[creds.Username]
    if !roleExists {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "role not assigned"})
        return
    }

    // Генерация токена с ролью пользователя
    token, err := generateToken(creds.Username, role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            if err == jwt.ErrSignatureInvalid {
                c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
                c.Abort() // Прерываем обработку запроса
                return
            }

            // Обработка истёкшего токена
            if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
                c.JSON(http.StatusUnauthorized, gin.H{"message": "token expired"})
                c.Abort()
                return
            }

            c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
            c.Abort()
            return
        }

        c.Next() // Если всё в порядке, передаём управление следующему обработчику
    }
}



var users = map[string]string{
    "admin": "admin123",
    "user": "password",
}

func register(c *gin.Context) {
    var creds Credentials
    if err := c.BindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
        return
    }

    // Проверка, существует ли пользователь
    if _, exists := users[creds.Username]; exists {
        handleError(c, http.StatusConflict, "user already exists")
        
        return
    }

    // По умолчанию роль "user", можно добавить проверку или параметр для роли
    role := "user" // Устанавливаем роль по умолчанию как "user"

    // Можно добавить параметр для роли в запросе регистрации, например:
    if creds.Role != "" {
        role = creds.Role
    }

    // Регистрируем пользователя
    users[creds.Username] = creds.Password
    roles[creds.Username] = role // Сохраняем роль в мапе

    c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}


var roles = map[string]string{
    "admin": "admin",
    "user": "user",
}

func roleMiddleware(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            handleError(c, http.StatusUnauthorized, "unauthorized")

            c.Abort()
            return
        }

        // Проверяем роль пользователя
        if claims.Role != requiredRole {
            handleError(c, http.StatusForbidden, "forbidden")

            c.Abort()
            return
        }

        c.Next()
    }
}


func refresh(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")
    claims := &Claims{}
    
    // Парсим исходный токен
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        handleError(c, http.StatusUnauthorized, "unauthorized")

        return
    }

    // Проверяем, не истек ли срок действия токена
    if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
        handleError(c, http.StatusBadRequest, "token not expired enough")
        return
    }

    // Генерация нового токена с теми же данными (пользователь и роль), но с новым временем истечения
    newToken, err := generateToken(claims.Username, claims.Role)
    if err != nil {
        handleError(c, http.StatusInternalServerError, "token not create token")

        return
    }

    c.JSON(http.StatusOK, gin.H{"token": newToken})
}

var db *gorm.DB

func initDB() {
    dsn := "host=localhost user=postgres password=67 dbname=test_store port=5432 sslmode=disable"
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    db.AutoMigrate(&Product{}, &Category{})
}




type Category struct {
    ID       uint      `gorm:"primaryKey" json:"id"`
    Name     string    `json:"name"`
    Products []Product `gorm:"foreignKey:CategoryID" json:"products"` // Связь с продуктами
}

type Product struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    CategoryID  uint   `json:"category_id"`
    Price       float64 `json:"price"`
    Manufacturer string `json:"manufacturer"`
}

func handleError(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, gin.H{"error": message})
}


func main() {
    initDB()
	router := gin.Default()

    router.POST("/login", login)
	router.POST("/register", register)
	router.POST("/refresh", refresh)

    protected := router.Group("/")
    protected.Use(authMiddleware())
    {
        protected.GET("/products/price-range", getProductsByPriceRange)
        protected.PUT("/products/manufacturer", roleMiddleware("admin"), updateProductsManufacturer)
        protected.GET("/products/count-by-manufacturer", countProductsByManufacturer)

    	protected.GET("/products", getProductsWithTimeout)

		protected.GET("/products/:id", getProductByID)

		protected.POST("/products", roleMiddleware("admin"), createProduct)

		protected.PUT("/products/:id", roleMiddleware("admin"), updateProduct)

		protected.DELETE("/products/:id", roleMiddleware("admin"), deleteProduct)

        protected.GET("/categories", getCategoriesWithTimeout)        // Получение всех категорий
		protected.GET("/categories/:id", getCategoryByID)      // Получение категории по ID
		protected.POST("/categories", roleMiddleware("admin"), createCategory)         // Создание новой категории
		protected.PUT("/categories/:id", roleMiddleware("admin"), updateCategory)       // Обновление категории
		protected.DELETE("/categories/:id", roleMiddleware("admin"), deleteCategory)    // Удаление категории

        
    }



	

	router.Run(":8080")
}


func getProducts(c *gin.Context) {
    var products []Product
    var total int64

    // Получаем параметры фильтров, сортировки и пагинации
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "10")
    sort := c.DefaultQuery("sort", "id")
    order := c.DefaultQuery("order", "asc")
    name := c.Query("name")
    categoryID := c.Query("category_id")

    // Преобразуем строковые параметры в int
    pageInt, _ := strconv.Atoi(page)
    limitInt, _ := strconv.Atoi(limit)
    offset := (pageInt - 1) * limitInt

    query := db.Model(&Product{})

    // Применяем фильтры
    if name != "" {
        query = query.Where("name ILIKE ?", "%"+name+"%")
    }
    if categoryID != "" {
        query = query.Where("category_id = ?", categoryID)
    }

    query.Count(&total)

    // Применяем сортировку
    if order != "asc" && order != "desc" {
        order = "asc" // По умолчанию ascending
    }
    query = query.Order(sort + " " + order).Limit(limitInt).Offset(offset)

    // Загружаем продукты и считаем общее количество
    query.Find(&products) 

    // Возращаем результат
    c.JSON(http.StatusOK, gin.H{
        "data":  products,
        "total": total,
        "page":  pageInt,
        "limit": limitInt,
    })
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")
    var product Product
    if err := db.First(&product, id).Error; err != nil {
        handleError(c, http.StatusNotFound, "Product not found")
        return
    }
    c.JSON(http.StatusOK, product)

}

func createProduct(c *gin.Context) {
	var newProduct Product

	if err := c.BindJSON(&newProduct); err != nil {
        handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

    var category Category
    if err := db.First(&category, newProduct.CategoryID).Error; err != nil {
        handleError(c, http.StatusBadRequest, "Invalid category ID")
        return
    }

    if newProduct.Price <= 0 {
        handleError(c, http.StatusBadRequest, "Price must be greater than 0")
        return
    }

	db.Create(&newProduct)
	c.JSON(http.StatusCreated, newProduct)

}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct Product

	if err := c.BindJSON(&updatedProduct); err != nil {
        handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

    if updatedProduct.Price <= 0 {
        handleError(c, http.StatusBadRequest, "Price must be greater than 0")
        return
    }

    if err := db.Model(&Product{}).Where("id = ?", id).Updates(updatedProduct).Error; err != nil {
        handleError(c, http.StatusNotFound, "Product not found")
        return
    }

	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")

    if err := db.Delete(&Product{}, id).Error; err != nil {
        handleError(c, http.StatusNotFound, "Product not found")
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})

	
}

// Получение всех категорий
func getCategories(c *gin.Context) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
        handleError(c, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}
	c.JSON(http.StatusOK, categories)
}

// Получение категории по ID
func getCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category Category
	if err := db.First(&category, id).Error; err != nil {
        handleError(c, http.StatusNotFound, "Category not found")
		return
	}
	c.JSON(http.StatusOK, category)
}

// Создание новой категории
func createCategory(c *gin.Context) {
	var newCategory Category
	if err := c.BindJSON(&newCategory); err != nil {
        handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := db.Create(&newCategory).Error; err != nil {
        handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}
	c.JSON(http.StatusCreated, newCategory)
}

// Обновление категории по ID
func updateCategory(c *gin.Context) {
	id := c.Param("id")
	var updatedCategory Category
	if err := c.BindJSON(&updatedCategory); err != nil {
        handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	// Проверяем, существует ли категория с этим ID
	var category Category
	if err := db.First(&category, id).Error; err != nil {
        handleError(c, http.StatusNotFound, "Category not found")
		return
	}

	// Обновляем категорию
	if err := db.Model(&category).Updates(updatedCategory).Error; err != nil {
        handleError(c, http.StatusInternalServerError, "Failed to update category")
		return
	}

    

	c.JSON(http.StatusOK, updatedCategory)
}

// Удаление категории по ID
func deleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Category{}, id).Error; err != nil {
        handleError(c, http.StatusNotFound, "Category not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}


func getProductsWithTimeout(c *gin.Context) {
    // Создаем контекст с тайм-аутом 2 секунды
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    var products []Product
    var total int64

    // Получаем параметры фильтров, сортировки и пагинации
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "10")
    sort := c.DefaultQuery("sort", "id")
    order := c.DefaultQuery("order", "asc")
    name := c.Query("name")
    categoryID := c.Query("category_id")

    // Преобразуем строковые параметры в int
    pageInt, _ := strconv.Atoi(page)
    limitInt, _ := strconv.Atoi(limit)
    offset := (pageInt - 1) * limitInt

    query := db.Model(&Product{})

    // Применяем фильтры
    if name != "" {
        query = query.Where("name ILIKE ?", "%"+name+"%")
    }
    if categoryID != "" {
        query = query.Where("category_id = ?", categoryID)
    }

    query.Count(&total)

    // Применяем сортировку
    if order != "asc" && order != "desc" {
        order = "asc" // По умолчанию ascending
    }
    query = query.Order(sort + " " + order).Limit(limitInt).Offset(offset)

    // Загружаем продукты с использованием контекста
    if err := query.WithContext(ctx).Find(&products).Error; err != nil {
        if err == context.DeadlineExceeded {
            handleError(c, http.StatusRequestTimeout, "Request timed out")
        } else {
            handleError(c, http.StatusInternalServerError, "Failed to fetch products")
        }
        return
    }

    // Возвращаем результат
    c.JSON(http.StatusOK, gin.H{
        "data":  products,
        "total": total,
        "page":  pageInt,
        "limit": limitInt,
    })
}

func getCategoriesWithTimeout(c *gin.Context) {
    // Создаем контекст с тайм-аутом 2 секунды
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    var categories []Category
    if err := db.WithContext(ctx).Preload("Products").Find(&categories).Error; err != nil {
        if err == context.DeadlineExceeded {
            handleError(c, http.StatusRequestTimeout, "Request timed out")
        } else {
            handleError(c, http.StatusInternalServerError, "Failed to fetch categories")
        }
        return
    }

    c.JSON(http.StatusOK, categories)
}

func getProductsByPriceRange(c *gin.Context) {
    minPrice, err1 := strconv.ParseFloat(c.Query("minPrice"), 64)
    maxPrice, err2 := strconv.ParseFloat(c.Query("maxPrice"), 64)
    
    if err1 != nil || err2 != nil {
        handleError(c, http.StatusBadRequest, "Invalid price range values")
        return
    }
    
    var products []Product
    if err := db.Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products).Error; err != nil {
        handleError(c, http.StatusInternalServerError, "Error fetching products")
        return
    }
    
    if len(products) == 0 {
        handleError(c, http.StatusNotFound, "No products found in the given price range")
        return
    }
    
    c.JSON(http.StatusOK, products)
    
}



func updateProductsManufacturer(c *gin.Context) {
    manufacturer := c.Query("manufacturer")

    if manufacturer == "" {
        handleError(c, http.StatusBadRequest, "manufacturer query parameter is required")
        return
    }

    // начало транзакции
    tx := db.Begin()

    // проверяем, что транзакция инициализирована корректно
    if tx.Error != nil {
        log.Println("Error starting transaction:", tx.Error)
        handleError(c, http.StatusInternalServerError, "Failed to start transaction")
        return
    }
    log.Println("Transaction started successfully.")

    // попытка массового обновления
    if err := tx.Model(&Product{}).Where("1 = 1").Update("manufacturer", manufacturer).Error; err != nil {
        tx.Rollback() // откатываем изменения при ошибке
        log.Println("Error during update operation:", err)
        handleError(c, http.StatusInternalServerError, "Error updating manufacturer: "+err.Error())
        return
    }
    log.Println("Manufacturer update operation successful.")

    // коммит транзакции
    if err := tx.Commit().Error; err != nil {
        log.Println("Error committing transaction:", err)
        handleError(c, http.StatusInternalServerError, "Transaction commit failed: "+err.Error())
        return
    }
    log.Println("Transaction committed successfully.")

    c.JSON(http.StatusOK, gin.H{"message": "Manufacturer updated successfully"})
}


func countProductsByManufacturer(c *gin.Context) {
    var result []struct {
        Manufacturer string
        Count        int
    }

    // Выполняем агрегацию по производителю и подсчитываем количество товаров
    if err := db.Model(&Product{}).
        Select("manufacturer, COUNT(*) as count").
        Group("manufacturer").
        Scan(&result).Error; err != nil {
        // Обработка ошибки, если что-то пошло не так
        handleError(c, http.StatusInternalServerError, "Error counting products by manufacturer: "+err.Error())
        return
    }

    // Возвращаем результат
    c.JSON(http.StatusOK, result)
}


```


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
