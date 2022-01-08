FROM golang:1.16.5 as golang-dev

LABEL maintainer=""

# agrego el directorio de trabajo
WORKDIR /app
# copiar las dependencias
COPY go.mod go.sum ./

RUN go mod download
# copiar todos los archivos
COPY . .
# installa reflex que sirve para ver los cambios en ejecucion 
RUN go install github.com/cespare/reflex@latest

# Expose para el puerto que se visualiza la aplicacion
EXPOSE 4000
# inicia app
CMD reflex -g '*.go' go run main.go --start-service
