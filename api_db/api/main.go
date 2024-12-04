package main


import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
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
        c.JSON(http.StatusConflict, gin.H{"message": "user already exists"})
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
            c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
            c.Abort()
            return
        }

        // Проверяем роль пользователя
        if claims.Role != requiredRole {
            c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
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
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
        return
    }

    // Проверяем, не истек ли срок действия токена
    if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
        c.JSON(http.StatusBadRequest, gin.H{"message": "token not expired enough"})
        return
    }

    // Генерация нового токена с теми же данными (пользователь и роль), но с новым временем истечения
    newToken, err := generateToken(claims.Username, claims.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
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
    //Products []Product `gorm:"foreignKey:CategoryID"` // Связь с продуктами
}

type Product struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    CategoryID  uint   `json:"category"`
    //Category    Category `gorm:"foreignKey:CategoryID"`
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
    	protected.GET("/products", getProducts)

		protected.GET("/products/:id", getProductByID)

		protected.POST("/products", roleMiddleware("admin"), createProduct)

		protected.PUT("/products/:id", roleMiddleware("admin"), updateProduct)

		protected.DELETE("/products/:id", roleMiddleware("admin"), deleteProduct)

        protected.GET("/categories", getCategories)        // Получение всех категорий
		protected.GET("/categories/:id", getCategoryByID)      // Получение категории по ID
		protected.POST("/categories", roleMiddleware("admin"), createCategory)         // Создание новой категории
		protected.PUT("/categories/:id", roleMiddleware("admin"), updateCategory)       // Обновление категории
		protected.DELETE("/categories/:id", roleMiddleware("admin"), deleteCategory)    // Удаление категории
    }



	

	router.Run(":8080")
}



/*



func getProducts(c *gin.Context) {
    var products []Product
    if err := db.Preload("Category").Find(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "error retrieving products"})
        return
    }
    c.JSON(http.StatusOK, products)
}
*/

func getProducts(c *gin.Context) {
    var products []Product
    db.Find(&products)
	c.JSON(http.StatusOK, products)
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")
    var product Product
    if err := db.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
    }
    c.JSON(http.StatusOK, product)

}

func createProduct(c *gin.Context) {
	var newProduct Product

	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

    var category Category
    if err := db.First(&category, newProduct.CategoryID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid category ID"})
        return
    }

	db.Create(&newProduct)
	c.JSON(http.StatusCreated, newProduct)

}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct Product

	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

    if err := db.Model(&Product{}).Where("id = ?", id).Updates(updatedProduct).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
    }

	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")

    if err := db.Delete(&Product{}, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})

	
}

// Получение всех категорий
func getCategories(c *gin.Context) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// Получение категории по ID
func getCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// Создание новой категории
func createCategory(c *gin.Context) {
	var newCategory Category
	if err := c.BindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if err := db.Create(&newCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, newCategory)
}

// Обновление категории по ID
func updateCategory(c *gin.Context) {
	id := c.Param("id")
	var updatedCategory Category
	if err := c.BindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	// Проверяем, существует ли категория с этим ID
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return
	}

	// Обновляем категорию
	if err := db.Model(&category).Updates(updatedCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// Удаление категории по ID
func deleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Category{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

