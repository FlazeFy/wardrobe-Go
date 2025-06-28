package utils

import (
	"fmt"
	"time"
	"wardrobe/models"

	"github.com/jung-kurt/gofpdf"
)

func GetStringNoData(s *string) string {
	if s != nil {
		return *s
	}
	return "-"
}

func GeneratePDFCreateClothes(c *models.Clothes, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Wardrobe", false)
	pdf.AddPage()

	// Set Header
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 102, 204)
	pdf.CellFormat(0, 12, "Wardrobe", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "Effortless style decision and Organize", "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// Set Letterhead
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("clothes : %s", c.ClothesName))
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("ID : %s", c.ID.String()))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Category : %s | Type : %s", c.ClothesCategory, c.ClothesType))
	pdf.Ln(10)

	// Set Props
	pdf.Cell(0, 10, fmt.Sprintf("Generated at: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(12)

	// Clothes Detail
	pdf.SetFont("Arial", "", 10)
	tableData := []struct {
		Label string
		Value string
	}{
		{"Name", c.ClothesName},
		{"Category", c.ClothesCategory},
		{"Description", GetStringNoData(c.ClothesDesc)},
		{"Merk", GetStringNoData(c.ClothesMerk)},
		{"Color", c.ClothesColor},
		{"Price", func() string {
			if c.ClothesPrice != nil {
				return fmt.Sprintf("Rp. %d", *c.ClothesPrice)
			}
			return "-"
		}()},
		{"Size", c.ClothesSize},
		{"Gender", c.ClothesGender},
		{"Made From", c.ClothesMadeFrom},
		{"Type", c.ClothesType},
		{"Purchased At", func() string {
			if c.ClothesBuyAt != nil {
				return c.ClothesBuyAt.Format("2006-01-02")
			}
			return "-"
		}()},
		{"Quantity", fmt.Sprintf("%d", c.ClothesQty)},
		{"Is Faded", BoolToYesNo(c.IsFaded)},
		{"Has Been Washed", BoolToYesNo(c.HasWashed)},
		{"Has Been Ironed", BoolToYesNo(c.HasIroned)},
		{"Is Favorite", BoolToYesNo(c.IsFavorite)},
		{"Is Scheduled", BoolToYesNo(c.IsScheduled)},
	}

	// Table Format
	pdf.SetFillColor(240, 240, 240)
	for _, row := range tableData {
		pdf.CellFormat(50, 8, row.Label, "1", 0, "C", true, 0, "")
		pdf.CellFormat(130, 8, row.Value, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	return pdf.OutputFileAndClose(filename)
}

func GeneratePDFErrorAudit(c []models.ErrorAudit, filename string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetTitle("Wardrobe", false)
	pdf.AddPage()

	// Set Header
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 102, 204)
	pdf.CellFormat(0, 12, "Wardrobe", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "Effortless style decision and Organize", "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// Set Letterhead
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Audit - Error")
	pdf.Ln(8)

	// Set header
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(175, 10, "Error Message", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 10, "Datetime", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Total", "1", 1, "C", true, 0, "")

	// Set body
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	for _, dt := range c {
		pdf.CellFormat(175, 8, dt.Message, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 8, dt.CreatedAt, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", dt.Total), "1", 1, "C", false, 0, "")
	}

	return pdf.OutputFileAndClose(filename)
}

func GeneratePDFReminderUnansweredQuestion(c []models.UnansweredQuestion, filename string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetTitle("Wardrobe", false)
	pdf.AddPage()

	// Set Header
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 102, 204)
	pdf.CellFormat(0, 12, "Wardrobe", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 10, "Effortless style decision and Organize", "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// Set Letterhead
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Reminder - Unanswered Question")
	pdf.Ln(8)

	// Set header
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(40, 10, "Date", "1", 0, "C", true, 0, "")
	pdf.CellFormat(150, 10, "Question", "1", 0, "C", true, 0, "")
	pdf.CellFormat(150, 10, "Answer", "1", 1, "C", true, 0, "")

	// Set body
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	for _, dt := range c {
		pdf.CellFormat(40, 8, dt.CreatedAt.Format("2006-01-02 15:04:05"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(150, 8, dt.Question, "1", 0, "L", false, 0, "")
		pdf.CellFormat(150, 8, "", "1", 1, "C", false, 0, "")
	}

	return pdf.OutputFileAndClose(filename)
}
