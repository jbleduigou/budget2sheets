package reader

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestReadUsingJsonBody(t *testing.T) {
	r := NewReader()

	m := events.SQSMessage{
		MessageId: "ID",
		Body: `
		{
			"Date": "28/02/2020", 
			"Description": "Brulerie Des Capuci Brest 28/02 Paiement Par Carte",
			"Comment": "Du bon café !",
			"Category": "Courses Alimentation",
			"Value": 6.25
		}
		`,
	}

	result, _ := r.Read(m)

	assert.Equal(t, result.Date, "28/02/2020")
	assert.Equal(t, result.Description, "Brulerie Des Capuci Brest 28/02 Paiement Par Carte")
	assert.Equal(t, result.Comment, "Du bon café !")
	assert.Equal(t, result.Category, "Courses Alimentation")
	assert.Equal(t, result.Value, 6.25)

}

func TestReadUsingJsonBodyInvalidFormat(t *testing.T) {
	r := NewReader()

	m := events.SQSMessage{
		MessageId: "ID",
		Body: `
		{
			"Date": "28/02/2020", 
			"Description": "Brulerie Des Capuci Brest 28/02 Paiement Par Carte",
			"Comment": "Du bon café !",
			"Category": "Courses Alimentation",
			"Value": "6.25"
		}
		`,
	}

	result, err := r.Read(m)

	assert.Equal(t, "json: cannot unmarshal string into Go struct field Transaction.Value of type float64", err.Error())

	assert.Equal(t, "28/02/2020", result.Date)
	assert.Equal(t, "Brulerie Des Capuci Brest 28/02 Paiement Par Carte", result.Description)
	assert.Equal(t, "Du bon café !", result.Comment)
	assert.Equal(t, "Courses Alimentation", result.Category)
	assert.Equal(t, 0.0, result.Value)
}
