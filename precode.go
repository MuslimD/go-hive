package main

import (
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"
)

var cafeList = map[string][]string{
    "moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
    countStr := req.URL.Query().Get("count")
    if countStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("count missing"))
        return
    }

    count, err := strconv.Atoi(countStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong count value"))
        return
    }

    city := req.URL.Query().Get("city")

    cafe, ok := cafeList[city]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong city value"))
        return
    }

    if count > len(cafe) {
        count = len(cafe)
    }

    answer := strings.Join(cafe[:count], ",")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    totalCount := 4
    req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)
    currentTotal := len(strings.Split(responseRecorder.Body.String(), ","))
    if currentTotal != totalCount {
        t.Errorf("expected total %d, got %d", totalCount, currentTotal)
    }
}

func TestMainHandlerWhenThereIsNoCity(t *testing.T) {
    req := httptest.NewRequest("GET", "/cafe?count=4&city=berlin", nil)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)
    body := responseRecorder.Body.String()
    if body != "wrong city value" {
        t.Errorf("expected total 'wrong city value', got %s", body)
    }
    code := responseRecorder.Code
    if code != 400 {
        t.Errorf("expected code 400, got %d", code)
    }
}

func TestMainHandlerWhenSuccess(t *testing.T) {
    req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)
    body := responseRecorder.Body
    if body == nil {
        t.Errorf("expected body, got nil")
    }
    code := responseRecorder.Code
    if code != 200 {
        t.Errorf("expected code 200, got %d", code)
    }
}

