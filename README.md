Markdown

### API GEMINI
API RESTful para gestión de usuarios y para procesar solicitudes de la API de Gemini de forma asíncrona.

---

### 🚀 **Cómo ejecutar el proyecto**

Sigue estos pasos para poner en marcha la aplicación en tu entorno local.

#### 1. Clonar el repositorio

```bash
git clone [https://github.com/Efren-Garza-Z/go-api-service.git](https://github.com/Efren-Garza-Z/go-api-gemini.git)
cd go-api-service
```

### 2. Configurar el entorno
Crea un archivo .env en la raíz del proyecto y agrega tus credenciales para la base de datos PostgreSQL y la clave de la API de Gemini.

Ini, TOML

DB_HOST=localhost
DB_PORT=5432
DB_USER=edgz
DB_PASSWORD=1234
DB_NAME=edgz
GEMINI_API_KEY=TU_API_KEY_DE_GEMINI

### 3. Instalar dependencias
Asegúrate de tener Go instalado. Luego, ejecuta el siguiente comando para instalar las dependencias del proyecto:

```bash

go mod tidy

```
### 4. Correr la aplicación
Puedes iniciar la API con el siguiente comando:

```bash

go run main.go
```

La aplicación se ejecutará en http://localhost:8080.

📖 Documentación de la API (Swagger)
La API utiliza Swagger para generar documentación interactiva.

1. Instalar la herramienta swag
Si no tienes swag instalado, debes agregarlo a tu sistema con el siguiente comando:

```bash

go install [github.com/swaggo/swag/cmd/swag@latest](https://github.com/swaggo/swag/cmd/swag@latest)
```
Nota para usuarios de Linux (Pop!_OS): El comando swag podría no ser reconocido directamente si no está en tu PATH. Puedes ejecutarlo con la ruta completa, que generalmente se encuentra en el directorio go/bin. Para verificar la ruta, usa el comando go env GOPATH.

2. Generar la documentación
Para generar o actualizar la documentación de Swagger, ejecuta el siguiente comando desde la terminal de tu proyecto:

```bash

$(go env GOPATH)/bin/swag init
```
Este comando buscará las anotaciones en tu código y creará o actualizará el archivo docs/swagger.json, que se utiliza para servir la documentación.

3. Acceder a la documentación
Una vez que la aplicación esté en funcionamiento, puedes acceder a la documentación interactiva en tu navegador:

http://localhost:8080/swagger/index.html