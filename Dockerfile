# Etapa de construcción
FROM golang:1.23.2 AS builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos de la aplicación Go al contenedor
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar la aplicación en modo de producción
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o urlshortener

# Etapa de producción
FROM alpine:latest

# Instalar las dependencias necesarias
RUN apk --no-cache add ca-certificates

# Copiar el binario compilado desde la etapa de construcción
COPY --from=builder /app/urlshortener /urlshortener

# Definir el puerto de la aplicación
ENV PORT=8080
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["/urlshortener"]