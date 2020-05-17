package reader

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var date = "<date/>"
var description = "<description/>"
var comment = "<comment/>"
var category = "<category/>"
var value = "-13.37"
var invalidValue = "trash"

func TestReadAttributesBeingPopulated(t *testing.T) {
	var messagetests = []struct {
		sqsDate        *string
		sqsDescription *string
		sqsComment     *string
		sqsCategory    *string
		sqsValue       *string
		date           string
		description    string
		comment        string
		category       string
		value          float64
	}{
		{nil, nil, nil, nil, nil, "", "", "", "", 0.0},
		{nil, nil, nil, nil, &invalidValue, "", "", "", "", 0.0},
		{nil, &description, &comment, &category, &value, "", description, comment, category, -13.37},
		{&date, nil, &comment, &category, &value, date, "", comment, category, -13.37},
		{&date, &description, nil, &category, &value, date, description, "", category, -13.37},
		{&date, &description, &comment, nil, &value, date, description, comment, "", -13.37},
		{&date, &description, &comment, &category, nil, date, description, comment, category, 0.0},
		{&date, &description, &comment, &category, &value, date, description, comment, category, -13.37},
	}

	r := NewReader()

	for _, v := range messagetests {
		m := events.SQSMessage{
			MessageId:         "ID",
			MessageAttributes: make(map[string]events.SQSMessageAttribute),
		}
		m.MessageAttributes["Date"] = events.SQSMessageAttribute{StringValue: v.sqsDate}
		m.MessageAttributes["Description"] = events.SQSMessageAttribute{StringValue: v.sqsDescription}
		m.MessageAttributes["Comment"] = events.SQSMessageAttribute{StringValue: v.sqsComment}
		m.MessageAttributes["Category"] = events.SQSMessageAttribute{StringValue: v.sqsCategory}
		m.MessageAttributes["Value"] = events.SQSMessageAttribute{StringValue: v.sqsValue}

		result, _ := r.Read(m)

		assert.Equal(t, result.Date, v.date)
		assert.Equal(t, result.Description, v.description)
		assert.Equal(t, result.Comment, v.comment)
		assert.Equal(t, result.Category, v.category)
		assert.Equal(t, result.Value, v.value)
	}

}

func TestReadAttributesNotBeingPopulated(t *testing.T) {
	r := NewReader()

	m := events.SQSMessage{
		MessageId:         "ID",
		MessageAttributes: make(map[string]events.SQSMessageAttribute),
	}

	result, _ := r.Read(m)

	assert.Equal(t, result.Date, "")
	assert.Equal(t, result.Description, "")
	assert.Equal(t, result.Comment, "")
	assert.Equal(t, result.Category, "")
	assert.Equal(t, result.Value, 0.0)

}

func TestReadUsingJsonBody(t *testing.T) {
	r := NewReader()

	m := events.SQSMessage{
		MessageId:         "ID",
		MessageAttributes: make(map[string]events.SQSMessageAttribute),
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

//For now we will accept this behaviour.
//When we remove support of Message Attributes we should return error instead.
//We might want as well to provide better support for incorrect SQS messages.
//This could be in the form of Dead Letter Queues and/or Cloud Watch Metrics and Alerts.
func TestReadUsingJsonBodyInvalidFormat(t *testing.T) {
	r := NewReader()

	m := events.SQSMessage{
		MessageId:         "ID",
		MessageAttributes: make(map[string]events.SQSMessageAttribute),
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

	result, _ := r.Read(m)

	assert.Equal(t, result.Date, "")
	assert.Equal(t, result.Description, "")
	assert.Equal(t, result.Comment, "")
	assert.Equal(t, result.Category, "")
	assert.Equal(t, result.Value, 0.0)
}
