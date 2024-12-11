# Imagen base para compilar aplicaciones Go
FROM golang:1.23-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos del proyecto al contenedor
COPY . .

# Compilar la aplicación
RUN go build -o main ./app/api/

# Imagen final para ejecutar la aplicación
FROM alpine:latest

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Expone el puerto de la aplicación
EXPOSE 8080

# Comando de entrada para ejecutar la aplicación
ENTRYPOINT ["./main"]
