package main

import (
	"fmt"
	"log"
	"net/http"
)

// 기본 핸들러: "/" 경로로 들어오는 요청을 처리
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// URL 경로가 정확히 "/" 가 아니라면 404 반환
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// 응답 헤더 설정
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// 본문에 문자열 출력
	fmt.Fprintln(w, "Welcome to Car Auction!")
}

func main() {
	// "/" 경로에 helloHandler 등록
	http.HandleFunc("/", helloHandler)

	// 8080 포트로 서버 시작
	addr := ":8080"
	log.Printf("Starting server on %s…\n", addr)
	// 에러 발생 시 로그 출력 후 종료
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
