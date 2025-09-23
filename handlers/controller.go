package handlers

import (
	"encoding/json"
	"errors"
	"go-rest-modul/database"
	"go-rest-modul/models"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"strconv"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ReadAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var recipes []models.Recipe

	err := database.DB.
		Preload("Category").
		Preload("RecipeIngredients").
		Preload("RecipeIngredients.Ingredient").
		Find(&recipes).Error

	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Fail to Query, Error : " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(recipes) == 0 {
		response := Response{
			Status:  "success",
			Message: "Recipe Not Found",
			Data:    recipes,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Status:  "success",
		Message: "All Recipe Retrieved Successfully",
		Data:    recipes,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ReadbyIDHandler(w http.ResponseWriter, r *http.Request) {
	var recipeid = mux.Vars(r)["id"]
	var recipe models.Recipe
	w.Header().Set("Content-Type", "application/json")
	err := database.DB.First(&recipe, recipeid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := Response{
				Status:  "error",
				Message: "Receip Not Found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		response := Response{
			Status:  "error",
			Message: "Error occured while retrieving data :" + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Status:  "success",
		Message: "Receipt Retrieved Successfully",
		Data:    recipe,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func AddRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var recipe models.Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Error occured while decoding data :" + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var categoryExist models.Category
	if err := database.DB.First(&categoryExist, recipe.CategoryId).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "Category Not Found",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = database.DB.Create(&recipe).Error
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Error occured while creating data :" + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := Response{
		Status:  "success",
		Message: "Receipt Created Successfully",
		Data:    recipe,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func UpdateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var recipeId = mux.Vars(r)["id"]
	var recipe models.Recipe

	if err := database.DB.First(&recipe, recipeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := Response{
				Status:  "error",
				Message: "Recipe Not Found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{
			Status:  "error",
			Message: "Database error: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Error occured while decoding data :" + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DB.Model(&recipe).Updates(recipe).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "Error occured while updating data :" + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Recipe Updated Successfully",
		Data:    recipe,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var recipeId = mux.Vars(r)["id"]
	var recipe models.Recipe
	w.Header().Set("Content-Type", "application/json")
	if err := database.DB.First(&recipe, recipeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := Response{
				Status:  "error",
				Message: "Recipe Not Found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{
			Status:  "error",
			Message: "Database error: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DB.Delete(&recipe).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "Error occured while deleting data :" + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Recipe Deleted Successfully",
		Data:    recipe,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SearchRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query().Get("q")
	var recipes []models.Recipe
	w.Header().Set("Content-Type", "application/json")
	if err := database.DB.
		Preload("Category").
		Preload("RecipeIngredients").
		Preload("RecipeIngredients.Ingredient").Preload("RecipeIngredients.Recipe").
		Where("title LIKE ?", "%"+query+"%").
		Or("descriptions LIKE ?", "%"+query+"%").
		Or("instructions LIKE ?", "%"+query+"%").
		Find(&recipes).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "Error occurred while searching data :" + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(recipes) == 0{
		response := Response{
			Status:  "error",
			Message: "Recipe Not Found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Search Succeed",
		Data:    recipes,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func FilterRecipesHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var recipes []models.Recipe

	w.Header().Set("Content-Type", "application/json")

	db := database.DB.Model(&models.Recipe{}).Preload("Category")

	if category := params.Get("category"); category != "" {
		db = db.Joins("JOIN categories ON categories.id = recipes.category_id").Where("categories.name = ?", category)
	}

	if maxpreptime := params.Get("max_preptime"); maxpreptime != "" {
		maxpreptime, err := strconv.Atoi(maxpreptime)
		if err == nil {
			db = db.Where("prep_time <= ?", maxpreptime)
		}
	}

	if servings := params.Get("servings"); servings != "" {
		servings, err := strconv.Atoi(servings)
		if err == nil {
			db = db.Where("servings = ?", servings)
		}
	}

	if err := db.Find(&recipes).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "Error occurred while filtering data :" + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(recipes) == 0 {
		response := Response{
			Status:  "error",
			Message: "Recipe Not Found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Filter Succeed",
		Data:    recipes,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func FilterByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var categoryId = mux.Vars(r)["category_id"]
	var recipes []models.Recipe

	w.Header().Set("Content-Type", "application/json")

	db := database.DB.Model(&models.Recipe{}).Preload("Category")

	if err := db.Find(&recipes, "category_id = ?", categoryId).Error; err != nil{
		response := Response{
			Status: "error",
			Message: "error occurred : " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(recipes) == 0{
		response := Response{
			Status: "error",
			Message: "Recipe Not Found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status: "success",
		Message: "Recipe by Category Found",
		Data: recipes,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}