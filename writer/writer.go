package writer

import (
	budget "github.com/jbleduigou/budget2sheets"
	"go.uber.org/zap"
	"google.golang.org/api/sheets/v4"
)

// Writer is a interface for writing a Transaction objects to Google Sheets
type Writer interface {
	Write(t budget.Transaction) error
}

// NewWriter returns an instances of a writer, actual implementation is not exposed
func NewWriter(srv *sheets.Service, spreadSheetID string, writeRange string, logger *zap.SugaredLogger) Writer {
	return &sheetsWriter{srv: srv, spreadSheetID: spreadSheetID, writeRange: writeRange, logger: logger}
}

type sheetsWriter struct {
	srv           *sheets.Service
	spreadSheetID string
	writeRange    string
	logger        *zap.SugaredLogger
}

func (w *sheetsWriter) Write(t budget.Transaction) error {
	w.logger.Infof("Processing SQS message with id '%v'", t.Date)
	vr, _ := asValueRange(t)
	_, err := w.srv.Spreadsheets.Values.Append(w.spreadSheetID, w.writeRange, &vr).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
	w.logger.Infof("Processed SQS message with id '%v'", t.Date)
	return err
}

func asValueRange(t budget.Transaction) (sheets.ValueRange, error) {
	var vr sheets.ValueRange
	myval := []interface{}{t.Date, t.Description, t.Comment, t.Category, t.Value}
	vr.Values = append(vr.Values, myval)
	return vr, nil
}
