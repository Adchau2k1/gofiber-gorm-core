# Sử dụng ảnh golang:latest làm cơ sở
FROM golang:latest

# Đặt thư mục làm việc
WORKDIR /app

# Sao chép mã nguồn vào container
COPY . .

# Tải dependencies và build ứng dụng
RUN go mod download
RUN go build -o app .

# CMD để chạy ứng dụng khi container được khởi chạy
CMD ["./app"]

