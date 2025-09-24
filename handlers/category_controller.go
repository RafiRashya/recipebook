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

func GetAllCategory(w http.ResponseWriter, r *http.Request){
	var category []models.Category

	w.Header().Set("Content-Type", "application/json")

	if err := database.DB.Preload("Recipes").Find(&category).Error; err != nil{
		response := Response{
			Status: "error",
			Message: "error occured while retrieving data: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(category) == 0{
		response := Response{
			Status: "not found",
			Message: "Category not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status : "success",
		Message: "Category retrieved successfully",
		Data: category,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetCategorybyId(w http.ResponseWriter, r *http.Request) {
	var categoryId = mux.Vars(r)["id"]
	var category models.Category

	w.Header().Set("Content-Type", "application/json")

	if err := database.DB.First(&category, categoryId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			response := Response{
				Status: "not found",
				Message: "Category not found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}
		response := Response{
			Status: "error",
			Message: "error occured: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status: "success",
		Message: "Category retrieved successfully",
		Data: category,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func CreateCategory(w http.ResponseWriter, r *http.Request){
	var category models.Category

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil{
		response := Response{
			Status: "error",
			Message: "Decode error :" + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DB.Create(&category).Error; err != nil {
		response := Response{
			Status: "error",
			Message: "Error while creating data : " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status: "success",
		Message: "Data created successfully",
		Data: category,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}