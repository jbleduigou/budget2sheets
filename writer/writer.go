package writer

import (
	"log/slog"

	budget "github.com/jbleduigou/budget2sheets"
	"google.golang.org/api/sheets/v4"
)

// Writer is a interface for writing a Transaction objects to Google Sheets
type Writer interface {
	Write(t budget.Transaction, messageId string) error
}

// NewWriter returns an instances of a writer, actual implementation is not exposed
func NewWriter(srv *sheets.Service, spreadSheetID string, writeRange string) Writer {
	return &sheetsWriter{srv: srv, spreadSheetID: spreadSheetID, writeRange: writeRange}
}

type sheetsWriter struct {
	srv           *sheets.Service
	spreadSheetID string
	writeRange    string
}

func (w *sheetsWriter) Write(t budget.Transaction, messageId string) error {
	slog.Info("Writing transaction to Google Sheets", slog.String("date", t.Date), slog.Float64("amount", t.Value))
	vr, _ := asValueRange(t)
	_, err := w.srv.Spreadsheets.Values.Append(w.spreadSheetID, w.writeRange, &vr).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
	slog.Info("Successfully wrote transaction to Google Sheets", slog.String("date", t.Date), slog.Float64("amount", t.Value))
	return err
}

func asValueRange(t budget.Transaction) (sheets.ValueRange, error) {
	var vr sheets.ValueRange
	myval := []interface{}{t.Date, t.Description, t.Comment, t.Category, t.Value}
	vr.Values = append(vr.Values, myval)
	return vr, nil
}
