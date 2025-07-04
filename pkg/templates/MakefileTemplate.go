package templates

const MakefileTemplate = `
# Generated by egg v0.0.1
# build-tailwindcss: # this is if you want to render your frontend
#   on the server without React 
# 	tailwindcss -i ./tailwind.css -o ./public/css/index.css 
build-backend:
	go build -o ./tmp/main .
build-frontend: 
	cd ./frontend/ && pnpm format && pnpm build
build-swagger:
	./tmp/main swag
build: build-backend build-swagger build-frontend
`


