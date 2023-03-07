package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_If_It_Get_An_Error_If_Did_Is_Blank(t *testing.T) {
	order := Order{}
	assert.Error(t, order.Validate(), "ID is required")
}

func Test_If_It_Get_An_Error_If_Price_Is_Less_Than_Zero(t *testing.T) {
	order := Order{ID: "123"}
	assert.Error(t, order.Validate(), "Price must be greater than zero")
}

func Test_If_It_Get_An_Error_If_Tax_Is_Less_Than_Zero(t *testing.T) {
	order := Order{ID: "123", Price: 10}
	assert.Error(t, order.Validate(), "Tax must be greater than zero")
}

func Test_If_It_Calculate_Final_Price(t *testing.T) {
	order := Order{ID: "123", Price: 10.0, Tax: 10.0}
	assert.NoError(t, order.CalculateFinalPrice())
	assert.Equal(t, 110.0, order.FinalPrice)
}

func Test_With_All_Valid_Params(t *testing.T) {
	order := Order{ID: "123", Price: 10.0, Tax: 10.0}
	assert.NoError(t, order.Validate())
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 10.0, order.Tax)
	
	order.CalculateFinalPrice()

	assert.Equal(t, 20.0, order.FinalPrice)
}
