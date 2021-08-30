touch .env
go get github.com/datadog/go-python3

archs=(amd64)

for arch in ${archs[@]}
do
        env GOOS=linux GOARCH=${arch} go build -o ./tmp/image_test_${arch} -i main/main.go
done