package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// UploadHandler は "file" フォームフィールドで送信されたファイルを受け取り
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストボディサイズを最大10MBに制限
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	// multipart フォームのパース
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// フォームから "file" フィールドを取得
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// uploads ディレクトリが存在しない場合は作成
	if err := os.MkdirAll("uploads", 0755); err != nil {
		http.Error(w, "could not create uploads dir: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 保存先ファイルを作成
	dstPath := filepath.Join("uploads", filepath.Base(header.Filename))
	out, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "could not save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// アップロードされたデータをコピー
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "could not write file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// アクセス用 URL を返却
	url := fmt.Sprintf("/static/%s", filepath.Base(header.Filename))
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"url": url}
	json.NewEncoder(w).Encode(resp)
}
