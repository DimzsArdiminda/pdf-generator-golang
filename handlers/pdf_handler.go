package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"pdf-api/config"
	"pdf-api/models"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
)
func DeletePDF(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	res, _ := config.DB.Exec(`
		UPDATE pdf_files
		SET status='DELETED', deleted_at=NOW()
		WHERE id=$1 AND status!='DELETED'
	`, id)

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "File not found or already deleted", 404)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF deleted successfully",
	})
}


func ListPDF(w http.ResponseWriter, r *http.Request) {
	rows, _ := config.DB.Query(`
		SELECT id, filename, original_name, size, status, created_at
		FROM pdf_files
		ORDER BY created_at DESC
	`)

	var data []models.PdfFile
	for rows.Next() {
		var p models.PdfFile
		rows.Scan(&p.ID, &p.Filename, &p.OriginalName, &p.Size, &p.Status, &p.CreatedAt)
		data = append(data, p)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File required", 400)
		return
	}
	defer file.Close()

	if filepath.Ext(header.Filename) != ".pdf" {
		http.Error(w, "Only PDF allowed", 400)
		return
	}

	filename := fmt.Sprintf("upload_%d.pdf", time.Now().Unix())
	path := "uploads/pdf/" + filename

	dst, _ := os.Create(path)
	defer dst.Close()
	io.Copy(dst, file)

	size := header.Size

	config.DB.Exec(`
		INSERT INTO pdf_files (filename, original_name, filepath, size, status)
		VALUES ($1,$2,$3,$4,'UPLOADED')
	`, filename, header.Filename, path, size)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF uploaded successfully",
	})
}

func GeneratePDF(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Title           string `json:"title"`
		InstitutionName string `json:"institution_name"`
		Address         string `json:"address"`
		Phone           string `json:"phone"`
		Content         string `json:"content"`
	}

	var req Req
	json.NewDecoder(r.Body).Decode(&req)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, req.InstitutionName)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, req.Address+" | "+req.Phone)
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, req.Title)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 8, req.Content, "", "", false)

	filename := fmt.Sprintf("report_%d.pdf", time.Now().Unix())
	path := "uploads/pdf/" + filename

	pdf.OutputFileAndClose(path)

	_, err := config.DB.Exec(`
		INSERT INTO pdf_files (filename, filepath, status)
		VALUES ($1, $2, 'CREATED')
	`, filename, path)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF generated successfully",
	})
}
