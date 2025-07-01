package templates

const DockerfileTemplate = `
## Build the Frontend with Node.js
## If you are not using React you can comment out this section
FROM node:22-alpine as node_builder
WORKDIR /usr/src/frontend
COPY frontend/package.json ./
## use pnpm
RUN npm install -g pnpm && pnpm install
COPY frontend/ ./
RUN pnpm run build

FROM golang:1.24-alpine as go_builder

WORKDIR /usr/src
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o {{.Name}} .

# Copy the executable to the final image
FROM alpine:latest as app

WORKDIR /app

## If you are not using React you can comment out this section
COPY --from=node_builder /usr/src/frontend/dist /app/{{.Server.Frontend.Dir}}

COPY --from=go_builder /usr/src/{{.Name}} /app/
CMD ["/app/{{.Name}}", "-e", "production"]
`

