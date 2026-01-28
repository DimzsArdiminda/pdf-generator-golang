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

	var deletedID int64
	var filename string
	var deletedAt time.Time
	
	err := config.DB.QueryRow(`
		UPDATE pdf_files
		SET status='DELETED', deleted_at=NOW()
		WHERE id=$1 AND status!='DELETED'
		RETURNING id, filename, deleted_at
	`, id).Scan(&deletedID, &filename, &deletedAt)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "File not found or already deleted",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF deleted successfully",
		"data": map[string]interface{}{
			"id":         deletedID,
			"filename":   filename,
			"status":     "DELETED",
			"deleted_at": deletedAt.Format(time.RFC3339),
		},
	})
}


func ListPDF(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := fmt.Sscanf(pageStr, "%d", &page); err == nil && p == 1 && page > 0 {
		} else {
			page = 1
		}
	}

	if limitStr != "" {
		if l, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && l == 1 && limit > 0 && limit <= 100 {
		} else {
			limit = 10
		}
	}

	offset := (page - 1) * limit

	var total int
	config.DB.QueryRow(`SELECT COUNT(*) FROM pdf_files WHERE status != 'DELETED'`).Scan(&total)

	rows, err := config.DB.Query(`
		SELECT id, filename, original_name, size, status, created_at
		FROM pdf_files
		WHERE status != 'DELETED'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var p models.PdfFile
		rows.Scan(&p.ID, &p.Filename, &p.OriginalName, &p.Size, &p.Status, &p.CreatedAt)
		
		item := map[string]interface{}{
			"id":         p.ID,
			"filename":   p.Filename,
			"size":       p.Size,
			"status":     p.Status,
			"created_at": p.CreatedAt.Format(time.RFC3339),
		}

		if p.OriginalName != nil {
			item["original_name"] = *p.OriginalName
		} else {
			item["original_name"] = nil
		}

		data = append(data, item)
	}

	
	if data == nil {
		data = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
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

	maxSize := int64(10 << 20)
	if header.Size > maxSize {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":    false,
			"message":    "File size exceeds maximum limit (10MB)",
			"error_code": "FILE_TOO_LARGE",
		})
		return
	}

	filename := fmt.Sprintf("upload_%d.pdf", time.Now().Unix())
	path := "uploads/pdf/" + filename

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "Failed to create file", 500)
		return
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", 500)
		return
	}

	size := header.Size

	var insertedID int64
	var createdAt time.Time
	err = config.DB.QueryRow(`
		INSERT INTO pdf_files (filename, original_name, filepath, size, status)
		VALUES ($1,$2,$3,$4,'UPLOADED')
		RETURNING id, created_at
	`, filename, header.Filename, path, size).Scan(&insertedID, &createdAt)

	if err != nil {
		http.Error(w, "Failed to save to database: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF uploaded successfully",
		"data": map[string]interface{}{
			"id":            insertedID,
			"original_name": header.Filename,
			"filename":      filename,
			"filepath":      path,
			"size":          size,
			"status":        "UPLOADED",
			"created_at":    createdAt.Format(time.RFC3339),
		},
	})
}

func GeneratePDF(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Title           string `json:"title"`
		InstitutionName string `json:"institution_name"`
		Address         string `json:"address"`
		Phone           string `json:"phone"`
		LogoURL         string `json:"logo_url"`
		Content         string `json:"content"`
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	if req.Title == "" || req.InstitutionName == "" || req.Content == "" {
		http.Error(w, "Title, institution_name, and content are required", 400)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	
	generatedTime := time.Now()
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb}", pdf.PageNo()),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Generated: %s", generatedTime.Format("02 Jan 2006 15:04")),
			"", 0, "R", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	})

	pdf.AliasNbPages("") 
	pdf.AddPage()

	pdf.SetFillColor(240, 240, 240)
	pdf.Rect(10, 10, 190, 35, "F")

	logoX := 15.0
	logoY := 15.0
	logoSize := 25.0
	
	if req.LogoURL != "" {
		pdf.Image(req.LogoURL, logoX, logoY, logoSize, 0, false, "", 0, "")
		if !pdf.Ok() {
			pdf.ClearError()
			pdf.SetDrawColor(200, 200, 200)
			pdf.Rect(logoX, logoY, logoSize, logoSize, "D")
			pdf.SetFont("Arial", "B", 8)
			pdf.SetXY(logoX+2, logoY+10)
			pdf.Cell(logoSize-4, 5, "LOGO")
		}
	} else {
		pdf.SetDrawColor(100, 100, 200)
		pdf.SetLineWidth(0.5)
		pdf.Rect(logoX, logoY, logoSize, logoSize, "D")
		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(logoX, logoY+10)
		pdf.CellFormat(logoSize, 5, "LOGO", "", 0, "C", false, 0, "")
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(45, 18)
	pdf.SetTextColor(0, 51, 102)
	pdf.Cell(0, 8, req.InstitutionName)

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(45, 28)
	pdf.SetTextColor(80, 80, 80)
	if req.Address != "" {
		pdf.Cell(0, 5, req.Address)
	}
	pdf.SetXY(45, 33)
	if req.Phone != "" {
		pdf.Cell(0, 5, "Tel: "+req.Phone)
	}

	pdf.SetTextColor(0, 0, 0)

	pdf.Ln(40)

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 51, 102)
	pdf.Cell(0, 10, req.Title)
	pdf.Ln(8)

	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.Cell(0, 6, "Date: "+generatedTime.Format("Monday, 02 January 2006"))
	pdf.Ln(10)

	pdf.SetDrawColor(0, 51, 102)
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(0, 0, 0)
	pdf.MultiCell(0, 6, req.Content, "", "J", false)

	filename := fmt.Sprintf("report_%d.pdf", time.Now().Unix())
	path := "uploads/pdf/" + filename

	if err := pdf.OutputFileAndClose(path); err != nil {
		http.Error(w, "Failed to generate PDF: "+err.Error(), 500)
		return
	}

	fileInfo, _ := os.Stat(path)
	fileSize := fileInfo.Size()

	var insertedID int64
	var createdAt time.Time
	err := config.DB.QueryRow(`
		INSERT INTO pdf_files (filename, original_name, filepath, size, status, logo_url)
		VALUES ($1, $2, $3, $4, 'CREATED', $5)
		RETURNING id, created_at
	`, filename, req.Title+".pdf", path, fileSize, req.LogoURL).Scan(&insertedID, &createdAt)

	if err != nil {
		http.Error(w, "Failed to save to database: "+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "PDF generated successfully",
		"data": map[string]interface{}{
			"id":         insertedID,
			"filename":   filename,
			"filepath":   path,
			"status":     "CREATED",
			"created_at": createdAt.Format(time.RFC3339),
		},
	})
}
