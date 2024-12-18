# Практическая работа 7, API (jwt)
## Шестериков Дмитрий ЭФМО-01-24
 
## 1. Регистрация нового пользователя
### Запрос
### Метод: POST
#### по умолчанию роль user
![](https://github.com/user-attachments/assets/05b05caf-ad9e-4d1d-960d-825cba085b97)

### Если вы попытаетесь зарегистрировать уже существующего пользователя, ответ будет следующим:
![image](https://github.com/user-attachments/assets/3b4abed8-0e11-46b9-8578-9c42add62e08)

## 2. Логин (получение токена)
### Запрос
### Метод: POST
![image](https://github.com/user-attachments/assets/bd79c031-db0a-45d1-adff-22b8e7b437b9)

### Если логин не прошел (например, неверный пароль), ответ будет следующим:
![image](https://github.com/user-attachments/assets/74416bd3-cfb3-491a-886b-e679dc054352)

## 3. Доступ к защищенным маршрутам
### Запрос
### Метод: GET
### Заголовки: В заголовке Authorization нужно передать действующий JWT токен
![image](https://github.com/user-attachments/assets/f3052498-4c98-4147-bd1f-219435919c9c)

## 4. Регистрация пользователя с ролью admin
### Запрос:
### Метод: POST
![image](https://github.com/user-attachments/assets/0e8a6e46-3663-41f4-852a-93e1640861ad)

## 5. Логин для получения токена (JWT) для admin
### Запрос:
### Метод: POST
![image](https://github.com/user-attachments/assets/830f96f6-cc6d-4f21-9b92-d7bf22840b27)

## 6. Доступ к защищенному ресурсу (только для роли admin)
### Теперь, когда у вас есть токен с ролью admin, вы можете получить доступ к защищенному ресурсу, который доступен только администраторам.
### Запрос:
### Метод: DELETE
![image](https://github.com/user-attachments/assets/e2d9a1ae-8733-4b24-9fda-127a61da501e)

## 7. Если вы попробуете доступ с ролью user
### Предположим, у вас есть пользователь с ролью user, и вы попытаетесь получить доступ к этому же защищенному ресурсу, который доступен только для администраторов.
### Запрос:
### Метод: DELETE
![image](https://github.com/user-attachments/assets/966acb04-b29a-4d02-83b2-3d7329c046c2)



```go
package main


import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
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




type Product struct {
	ID 			string
	Name		string
	Description string
	Category 	string
}

var products = []Product{
	{ID: "1", Name: "Протеин сывороточный", Description: "Высококачественный сывороточный протеин для роста мышц.", Category: "Протеины"},
	{ID: "2", Name: "Креатин моногидрат", Description: "Увеличивает силу и выносливость при тренировках.", Category: "Креатин"},
	{ID: "3", Name: "BCAA комплекс", Description: "Комплекс аминокислот для восстановления и роста мышц.", Category: "Аминокислоты"},
	{ID: "4", Name: "Витамины для спортсменов", Description: "Комплекс витаминов и минералов для поддержки здоровья.", Category: "Витамины"},
	{ID: "5", Name: "Гейнер", Description: "Высококалорийный продукт для быстрого набора массы.", Category: "Гейнеры"},
	{ID: "6", Name: "Омега-3", Description: "Полиненасыщенные жирные кислоты для здоровья сердца.", Category: "Добавки"},
	{ID: "7", Name: "Сжигатель жира", Description: "Продукт для контроля веса и ускорения метаболизма.", Category: "Сжигатели жира"},
	{ID: "8", Name: "Протеин растительный", Description: "Протеин на основе гороха и риса для вегетарианцев.", Category: "Протеины"},
	{ID: "9", Name: "Спортивные батончики", Description: "Батончики с высоким содержанием белка для перекуса.", Category: "Перекусы"},
	{ID: "10", Name: "Изотонический напиток", Description: "Увлажняющий напиток для восстановления во время тренировок.", Category: "Напитки"},
}

func main() {
	router := gin.Default()

    router.POST("/login", login)
	router.POST("/register", register)
	router.POST("/refresh", refresh)

    protected := router.Group("/")
    protected.Use(authMiddleware())
    {
    	router.GET("/products", getProducts)

		router.GET("/products/:id", getProductByID)

		router.POST("/products", roleMiddleware("admin"), createProduct)

		router.PUT("/products/:id", roleMiddleware("admin"), updateProduct)

		router.DELETE("/products/:id", roleMiddleware("admin"), deleteProduct)
    }



	

	router.Run(":8080")
}

func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")

	for _, product := range products {
		if product.ID == id {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
}

func createProduct(c *gin.Context) {
	var newProduct Product

	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	products = append(products, newProduct)
	c.JSON(http.StatusCreated, newProduct)

}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct Product

	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	for i, product := range products {
		if product.ID == id {
			products[i] = updatedProduct
			c.JSON(http.StatusOK, updatedProduct)
			return
		}
	}


}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i + 1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

```
