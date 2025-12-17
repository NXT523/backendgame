1. kiểm tra phiên bản
$env:Path += ";C:\tools\protoc\bin;C:\Users\$env:USERNAME\go\bin"
protoc --version
2. Cài plugin (chạy 1 lần trong terminal VS Code)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
-Nếu dùng REST gateway:
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
3. Tải file googleapis (nếu có REST)
cd C:\Golang\go_base\19_game (trỏ vào thư mục muốn lưu)
# 1) Tạo đúng thư mục (đã có cũng không sao)
New-Item -ItemType Directory -Force .\third_party\googleapis\google\api | Out-Null
# 2) Tải file về đúng đường dẫn
$BASE = "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api"
Invoke-WebRequest "$BASE/annotations.proto" -OutFile ".\third_party\googleapis\google\api\annotations.proto"
Invoke-WebRequest "$BASE/http.proto"         -OutFile ".\third_party\googleapis\google\api\http.proto"
# 3) Kiểm tra
Test-Path .\third_party\googleapis\google\api\annotations.proto
Test-Path .\third_party\googleapis\google\api\http.proto
4. Gọi protoc phải include đúng cả 2 thư mục
$env:Path += ";C:\tools\protoc\bin;C:\Users\$env:USERNAME\go\bin"
protoc -I proto -I third_party/googleapis `
  --go_out=. --go_opt=paths=source_relative `
  --go-grpc_out=. --go-grpc_opt=paths=source_relative `
  --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative `
  proto\v1\*.proto
5. 
go mod tidy