
generate-tailwind-css:
	npx tailwindcss -i templates/css/raw/input.css -o templates/css/my-tailwind.css --watch

build-dev:
	ENV=dev go build ./main.go -o web3-blog

run-dev:
	ENV=dev go run ./main.go

build-docker:
	docker build -t zou8944/web3-blog .

run-docker:
	docker run --name web3-blog -v ~/.web3-blog/config:/config -p 9000:9000 -e ENV=prod zou8944/web3-blog:latest