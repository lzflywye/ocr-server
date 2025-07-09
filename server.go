package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 64 << 20 // 64 MiB

	router.POST("/upload", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			data := map[string]interface{}{
				"success": false,
				"data":    map[string]interface{}{},
				"errors": map[string]interface{}{
					"code":    "INPUT ERROR",
					"message": "データの入力がありませんでした。",
				},
			}
			ctx.SecureJSON(http.StatusBadRequest, data)
		}

		ocrClient := gosseract.NewClient()
		ocrClient.SetLanguage("eng+jpn")
		defer ocrClient.Close()

		files := form.File["files"]
		for _, file := range files {
			src, err := file.Open()
			if err != nil {
				data := map[string]interface{}{
					"success": false,
					"data":    map[string]interface{}{},
					"errors": map[string]interface{}{
						"code":    "FILE OPEN ERROR",
						"message": "ファイルを開くことができませんでした",
					},
				}
				ctx.SecureJSON(http.StatusInternalServerError, data)
			}
			defer src.Close()

			fileBytes, err := io.ReadAll(src)
			if err != nil {
				data := map[string]interface{}{
					"success": false,
					"data":    map[string]interface{}{},
					"errors": map[string]interface{}{
						"code":    "FILE READ ERROR",
						"message": "ファイルを読み取ることができませんでした",
					},
				}
				ctx.SecureJSON(http.StatusInternalServerError, data)
			}

			err = ocrClient.SetImageFromBytes(fileBytes)
			if err != nil {
				data := map[string]interface{}{
					"success": false,
					"data":    map[string]interface{}{},
					"errors": map[string]interface{}{
						"code":    "OCR ERROR",
						"message": "OCRの実行に失敗しました",
					},
				}
				ctx.SecureJSON(http.StatusInternalServerError, data)
			}

			text, err := ocrClient.Text()
			if err != nil {
				data := map[string]interface{}{
					"success": false,
					"data":    map[string]interface{}{},
					"errors": map[string]interface{}{
						"code":    "OCR TEXT ERROR",
						"message": "OCRの文字列出力に失敗しました",
					},
				}
				ctx.SecureJSON(http.StatusInternalServerError, data)
			}

			data := map[string]interface{}{
				"success": true,
				"data": map[string]interface{}{
					"text": text,
				},
				"errors": map[string]interface{}{},
			}
			ctx.SecureJSON(http.StatusOK, data)
		}

	})

	router.Run("127.0.0.1:8080")
}
