.PHONY: tailwind-build
tailwind-build:
	nix-shell -p tailwindcss --run "tailwindcss -i ./assets/css/input.css -o ./assets/css/tailwind.css"
	nix-shell -p tailwindcss --run "tailwindcss -i ./assets/css/input.css -o ./assets/css/tailwind.min.css --minify"

.PHONY: tailwind-watch
tailwind-watch: 
	nix-shell -p tailwindcss --run "tailwindcss -i ./assets/css/input.css -o ./assets/css/tailwind.css --watch"

.PHONY: templ-generate
templ-generate:
	nix run github:a-h/templ generate

.PHONY: templ-watch
templ-watch:
	nix run github:a-h/templ -- generate --watch
	
.PHONY: dev
dev:
	make templ-watch &
	make tailwind-watch &	
	go build -o ./tmp/app ./main.go && air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o ./bin/app ./main.go


