# 🎬 CineConecta Backend

Este es el backend del proyecto **CineConecta**, una plataforma de gestión de películas desarrollada en Go usando el framework **Gin**. El backend proporciona autenticación segura mediante JWT, gestión de usuarios con roles y un CRUD completo para películas, incluyendo filtros dinámicos.

---

## 🚀 Tecnologías utilizadas

- **Go 1.20+**
- **Gin Gonic** (framework web)
- **GORM** (ORM para Go)
- **PostgreSQL** (Base de datos)
- **JWT** (Autenticación)
- **Vercel** (Despliegue en la nube)

---

## 🗂 Estructura del proyecto

```
├── auth/
│   ├── controllers/   # Lógica de los endpoints de autenticación
│   ├── middlewares/   # Validación de tokens y roles
│   ├── models/        # Modelo de Usuario
│   ├── services/      # Lógica de negocio del usuario
│   ├── utils/         # JWT, manejo de errores, cookies
│   └── routes/        # Rutas de autenticación
│
├── movies/
│   ├── controllers/   # Endpoints de películas
│   ├── models/        # Modelo de Película
│   ├── services/      # Lógica de negocio de películas
│   └── routes/        # Rutas de películas
│
├── config/            # Conexión a la base de datos
├── handler/           # Entry point para Vercel (index.go)
├── go.mod / go.sum    # Módulo de Go
└── vercel.json        # Configuración para despliegue
```

---

## 🔐 Autenticación y Roles

- Registro e inicio de sesión con contraseñas cifradas.
- Generación de token JWT que contiene el nombre, ID y rol del usuario.
- Token almacenado en cookie segura `cine_token`.
- Acciones restringidas a usuarios con rol `admin` (como crear/editar/eliminar películas).

---

## 🎞 Endpoints disponibles

### Usuarios
```
POST   /api/register         # Registro
POST   /api/login            # Login y seteo del token en cookie
POST   /api/logout           # Logout (elimina cookie)
GET    /api/profile          # Datos del usuario autenticado
GET    /api/users            # (admin) Ver todos los usuarios
DELETE /api/users            # (admin) Eliminar todos excepto admin
GET    /api/verify-token     # Verifica si el token es válido
```

### Películas
```
GET    /api/movies                  # Obtener todas
GET    /api/movies/sorted          # Obtener con ordenamiento dinámico (por ?sortBy=&order=)
GET    /api/movies/:id             # Obtener una por ID
POST   /api/movies                 # (admin) Crear nueva
PUT    /api/movies/:id            # (admin) Actualizar
DELETE /api/movies/:id            # (admin) Eliminar
```

---

## ⚙️ Filtro dinámico en películas

Puedes ordenar por `title`, `genre` o `rating`, por ejemplo:

```
GET /api/movies/sorted?sortBy=rating&order=desc
```

---

## 🌐 CORS habilitado

Permite peticiones desde:
- `http://localhost:3000`
- `https://tufrontend.vercel.app`

Con credenciales (cookies) activadas.

---

## 🛠 Configuración del entorno

Crea un archivo `.env` con:

```
# Configuración básica
DATABASE_URL=postgresql://... (tu cadena de conexión a PostgreSQL)
JWT_SECRET=clave-secreta-segura
ENV=development

# Configuración para almacenamiento S3 (para subida de imágenes)
S3_ENDPOINT=https://tu-proyecto.supabase.co/storage/v1/s3
S3_REGION=us-east-1
S3_BUCKET=nombre-de-tu-bucket
S3_ACCESS_KEY=tu-access-key
S3_SECRET_KEY=tu-secret-key
```

### Configuración de Supabase para almacenamiento

Para habilitar la subida de imágenes utilizando Supabase como proveedor S3:

1. Crea un bucket público en Supabase Storage llamado `imagenes` (o el nombre que prefieras)
2. Ve a la sección Storage > Policies de tu bucket y configura una política que permita uploads públicos
3. Copia las credenciales de la sección Project Settings > API a tu archivo `.env`

---

## 🚀 Despliegue en Vercel

Este backend está desplegado como **Serverless Function** en Vercel. El entry point es `handler/index.go`.

Archivos importantes:
- `vercel.json`: especifica rutas y runtime Go
- `go.mod`: define el módulo y dependencias

---

## 🧪 Datos de prueba
Puedes usar `seed_movies.go` para insertar películas de prueba en la base de datos sin duplicados.

---

## 👨‍💻 Autor

**Felipe Huertas**  
Backend Developer  
📧 fhuertas@unillanos.edu.co

**Juan Romero**
Backend Developer
📧 juanromero2719@gmail.com
https://wrydmoon.site
