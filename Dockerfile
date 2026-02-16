# 1. 選擇基底映像檔 (就像選擇作業系統)
# 我們使用官方的 Go 語言環境 (版本 1.25)
FROM golang:1.25

# 2. 設定工作目錄
# 接下來的指令都會在這個資料夾內執行 (類似 cd /app)
WORKDIR /app

# 3. 複製依賴檔案
# 先只複製 go.mod 和 go.sum，這是為了利用 Docker 的快取機制
COPY go.mod go.sum ./

# 4. 下載依賴
# 根據 go.mod 下載 Gin 等套件
RUN go mod download

# 5. 複製程式碼
# 把目前目錄下的所有東西 (main.go 等) 複製進容器的 /app
COPY . .

# 6. 編譯程式
# 把 main.go 編譯成一個執行檔，取名叫 main
RUN go build -o main .

# 7. 宣告端口
# 告訴外面的人，這個容器會用 8080 port
EXPOSE 8080

# 8. 啟動指令
# 當容器跑起來時，執行這個指令
CMD ["./main"]