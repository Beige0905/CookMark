package handler

import (
	"io"
	"net/http"

	"github.com/Beige0905/recipe-backend/internal/service"
)

type OCRHandler struct {
	ocrService *service.OCRService
}

func NewOCRHandler(os *service.OCRService) *OCRHandler {
	return &OCRHandler{
		ocrService: os,
	}
}

func (h *OCRHandler) ExtractImage(w http.ResponseWriter, r *http.Request) {
	// 1. 이미지 파일 추출
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		http.Error(w, "파일 파싱 실패", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "이미지 파일이 필요합니다", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "파일 읽기 실패", http.StatusInternalServerError)
		return
	}

	// 2. 서비스 호출
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	result, err := h.ocrService.ExtractFromImage(r.Context(), fileBytes, contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, result)
}
