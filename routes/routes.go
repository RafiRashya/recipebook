package routes

import (
	"github.com/gorilla/mux"
	"go-rest-modul/handlers"
)

func RegisterRoutes() *mux.Router{
	router := mux.NewRouter().StrictSlash(true)

	// Recipe Routes
	recipe := router.PathPrefix("/api/recipe").Subrouter()
	recipe.HandleFunc("", handlers.AddRecipeHandler).Methods("POST")
	recipe.HandleFunc("/{id}", handlers.ReadbyIDHandler).Methods("GET")
	recipe.HandleFunc("/{id}", handlers.UpdateRecipeHandler).Methods("PUT")
	recipe.HandleFunc("/{id}", handlers.DeleteRecipeHandler).Methods("DELETE")

	// Recipes Collection
	recipes := router.PathPrefix("/api/recipes").Subrouter()
	recipes.HandleFunc("", handlers.ReadAllHandler).Methods("GET")
	recipes.HandleFunc("/search", handlers.SearchRecipeHandler).Methods("GET")
	recipes.HandleFunc("/filter", handlers.FilterRecipesHandler).Methods("GET")
	recipes.HandleFunc("/category/{category_id}", handlers.FilterByCategoryHandler).Methods("GET")

	category := router.PathPrefix("/api/category").Subrouter()
	category.HandleFunc("/{{id}}", handlers.GetCategorybyId).Methods("GET")
	category.HandleFunc("", handlers.CreateCategory).Methods("POST")

	//routes untuk Category Functionality
	categories := router.PathPrefix("/api/categories").Subrouter()
	categories.HandleFunc("", handlers.GetAllCategory).Methods("GET")

	return router
}