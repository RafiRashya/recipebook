package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Title             string
	Descriptions      string
	Instructions      string
	PrepTime          int
	CookTime          int
	Servings          int
	ImageURL          string
	CategoryId        uint
	Category          Category           `gorm:"foreignKey:CategoryId"`
	RecipeIngredients []RecipeIngredient `gorm:"foreignKey:RecipeId"`
}
type Ingredient struct {
	gorm.Model
	Name              string
	RecipeIngredients []RecipeIngredient `gorm:"foreignKey:IngredientId"`
}

type Category struct {
	gorm.Model
	Name    string
	Recipes []Recipe `gorm:"foreignKey:CategoryId"`
}

type RecipeIngredient struct {
	gorm.Model
	RecipeId     uint
	Recipe       Recipe `gorm:"foreignKey:RecipeId"`
	IngredientId uint
	Ingredient   Ingredient `gorm:"foreignKey:IngredientId"`
	Amount       string
	Unit         string
}
