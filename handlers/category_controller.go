package handlers

import (
	"encoding/json"
	"errors"
	"go-rest-modul/database"
	"go-rest-modul/models"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	var category []models.Category

	w.Header().Set("Content-Type", "application/json")

	if err := database.DB.Preload("Recipes").Find(&category).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "error occured while retrieving data: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(category) == 0 {
		response := Response{
			Status:  "not found",
			Message: "Category not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Category retrieved successfully",
		Data:    category,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetCategorybyId(w http.ResponseWriter, r *http.Request) {
	var categoryId = mux.Vars(r)["id"]
	var category models.Category

	w.Header().Set("Content-Type", "application/json")

	if err := database.DB.First(&category, categoryId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := Response{
				Status:  "not found",
				Message: "Category not found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		response := Response{
			Status:  "error",
			Message: "error occured: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Category retrieved successfully",
		Data:    category,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		response := Response{
			Status:  "error",
			Message: "Decode error: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if category.Name == "" {
		response := Response{
			Status:  "error",
			Message: "Name cannot be empty",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var existing models.Category
	if err := database.DB.Where("name = ?", category.Name).First(&existing).Error; err == nil {
		response := Response{
			Status:  "error",
			Message: "Category already exists",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		response := Response{
			Status:  "error",
			Message: "Failed to check category: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DB.Create(&category).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "Error while creating data: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Data created successfully",
		Data:    category,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryId = mux.Vars(r)["id"]
	var category models.Category

	w.Header().Set("Content-Type", "application/json")

	// Cari data lama
	if err := database.DB.First(&category, categoryId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := Response{
				Status:  "error",
				Message: "Couldn't find category",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		response := Response{
			Status:  "error",
			Message: "Failed when checking data : " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Decode ke struct baru
	var input struct {
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response := Response{
			Status:  "error",
			Message: "error while decoding",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if input.Name == "" {
		response := Response{
			Status:  "error",
			Message: "Name cannot be empty",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Cek duplikasi nama (jika nama berubah)
	if input.Name != category.Name {
		var existing models.Category
		if err := database.DB.Where("name = ?", input.Name).First(&existing).Error; err == nil {
			response := Response{
				Status:  "error",
				Message: "Category name already exists",
			}
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			response := Response{
				Status:  "error",
				Message: "Failed to check category: " + err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Update field yang diizinkan
	category.Name = input.Name

	if err := database.DB.Save(&category).Error; err != nil {
		response := Response{
			Status:  "error",
			Message: "error when update: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Category has been updated",
		Data:    category,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var categoryId = mux.Vars(r)["id"]
	var category models.Category

	if err := database.DB.First(&category, categoryId).Error; err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			response := Response{
				Status: "error",
				Message: "Category Not Found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		response := Response{
			Status: "error",
			Message: "Error Occure while searching data : " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DB.Unscoped().Delete(&category, categoryId).Error; err != nil{
		response := Response{
			Status: "error",
			Message: "An error occured while deleting category" + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status: "success",
		Message: "Cateogry has been deleted successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
