package main


import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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

	router.GET("/products", getProducts)

	router.GET("/products/:id", getProductByID)

	router.POST("/products", createProduct)

	router.PUT("/products/:id", updateProduct)

	router.DELETE("/products/:id", deleteProduct)

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