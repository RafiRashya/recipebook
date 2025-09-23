package main

import (
	"github.com/gorilla/mux"
	"go-rest-modul/database"
	"go-rest-modul/handlers"
	"go-rest-modul/models"
	"log"
	"net/http"
)

func seedData() {
	// Seed Categories
	categories := []models.Category{
		{Name: "Main Course"},
		{Name: "Dessert"},
		{Name: "Breakfast"},
	}

	for _, category := range categories {
		database.DB.FirstOrCreate(&category, models.Category{Name: category.Name})
	}

	// Seed Ingredients
	ingredients := []models.Ingredient{
		{Name: "Bananas"},
		{Name: "Milk"},
		{Name: "Eggs"},
		{Name: "Sugar"},
		{Name: "Flour"},
		{Name: "Butter"},
	}

	for _, ingredient := range ingredients {
		database.DB.FirstOrCreate(&ingredient, models.Ingredient{Name: ingredient.Name})
	}

	// Get categories for foreign keys
	var breakfast, dessert models.Category
	database.DB.Where("name = ?", "Breakfast").First(&breakfast)
	database.DB.Where("name = ?", "Dessert").First(&dessert)

	// Seed Recipes
	recipes := []models.Recipe{
		{
			Title:        "Banana Pancakes",
			Descriptions: "Fluffy pancakes with fresh bananas",
			Instructions: "Mix ingredients and cook on griddle",
			PrepTime:     10,
			CookTime:     15,
			Servings:     4,
			ImageURL:     "https://example.com/pancakes.jpg",
			CategoryId:   breakfast.ID,
		},
		{
			Title:        "Sugar Cookies",
			Descriptions: "Sweet vanilla cookies",
			Instructions: "Mix, shape, and bake",
			PrepTime:     20,
			CookTime:     12,
			Servings:     24,
			ImageURL:     "https://example.com/cookies.jpg",
			CategoryId:   dessert.ID,
		},
	}

	for _, recipe := range recipes {
		database.DB.FirstOrCreate(&recipe, models.Recipe{Title: recipe.Title})
	}

	log.Println("Data berhasil di seed")
}

func main() {
	log.Println("Memulai seed data")
	seedData()
	//data, err := handlers.ReadAllHandler()
	//if err != nil {
	//	fmt.Println("ada error ", err.Error())
	//}
	//fmt.Println(data)

	log.Println("Selesai seed data")
	log.Println("Memulai server")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/recipes", handlers.ReadAllHandler).Methods("GET")
	router.HandleFunc("/api/recipe/{id}", handlers.ReadbyIDHandler).Methods("GET")
	router.HandleFunc("/api/recipe", handlers.AddRecipeHandler).Methods("POST")
	router.HandleFunc("/api/recipe/{id}", handlers.UpdateRecipeHandler).Methods("PUT")
	router.HandleFunc("/api/recipe/{id}", handlers.DeleteRecipeHandler).Methods("DELETE")
	router.HandleFunc("/api/recipes/search", handlers.SearchRecipeHandler).Methods("GET")
	router.HandleFunc("/api/recipes/filter", handlers.FilterRecipesHandler).Methods("GET")
	router.HandleFunc("/api/recipes/category/{category_id}", handlers.FilterByCategoryHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
