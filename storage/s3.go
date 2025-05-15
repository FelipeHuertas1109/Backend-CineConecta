// storage/s3.go
package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/joho/godotenv"
)

var (
	bucket     string
	baseURL    string
	uploader   *manager.Uploader
	isSupabase bool
)

func init() {
	// 1) Carga .env
	_ = godotenv.Load()

	// 2) Imprime configuración para depurar
	endpoint := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	region := strings.TrimSpace(os.Getenv("S3_REGION"))
	bucketEnv := strings.TrimSpace(os.Getenv("S3_BUCKET"))
	accessEnv := strings.TrimSpace(os.Getenv("S3_ACCESS_KEY"))
	secretEnv := strings.TrimSpace(os.Getenv("S3_SECRET_KEY"))

	fmt.Println("📦 [DEBUG] Configuración S3 cargada:")
	fmt.Println("    S3_ENDPOINT  :", endpoint)
	fmt.Println("    S3_REGION    :", region)
	fmt.Println("    S3_BUCKET    :", bucketEnv)
	fmt.Println("    S3_ACCESS_KEY:", maskSecret(accessEnv))
	fmt.Println("    S3_SECRET_KEY:", maskSecret(secretEnv))

	if endpoint == "" || region == "" || bucketEnv == "" || accessEnv == "" || secretEnv == "" {
		log.Println("⚠️ [ADVERTENCIA] Variables de entorno S3 incompletas. Subida de imágenes deshabilitada.")
		return
	}

	// 3) Asigna variables globales
	bucket = bucketEnv
	isSupabase = strings.Contains(endpoint, "supabase")

	// Para Supabase, la URL pública debe construirse correctamente
	if isSupabase {
		// Formato correcto para Supabase: https://[proyecto].supabase.co/storage/v1/object/public/[bucket]/[path]
		baseEndpoint := strings.Replace(endpoint, "/storage/v1/s3", "", 1)
		baseURL = fmt.Sprintf("%s/storage/v1/object/public/%s", baseEndpoint, bucket)
	} else {
		baseURL = fmt.Sprintf("%s/%s", endpoint, bucket)
	}

	fmt.Println("    BASE_URL     :", baseURL)
	fmt.Println("    ES SUPABASE  :", isSupabase)

	// 4) Configuración personalizada para Supabase
	var s3Client *s3.Client

	if isSupabase {
		// Para Supabase, usamos configuración específica
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, reg string, _ ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: reg,
			}, nil
		})

		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessEnv, secretEnv, "")),
			config.WithEndpointResolverWithOptions(customResolver),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Error al cargar configuración S3: %v", err))
		}

		// Cliente S3 para Supabase - sin path style
		s3Client = s3.NewFromConfig(cfg)
	} else {
		// Configuración estándar para otros proveedores S3
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessEnv, secretEnv, "")),
			config.WithEndpointResolverWithOptions(
				aws.EndpointResolverWithOptionsFunc(func(service, reg string, _ ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: reg,
						SigningName:   "s3",
					}, nil
				})),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Error al cargar configuración S3: %v", err))
		}

		// Cliente S3 estándar con path style
		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	// 6) Inicializa el uploader con timeouts
	uploader = manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.PartSize = 5 * 1024 * 1024 // 5MB por parte
		u.LeavePartsOnError = false
	})
}

// UploadPoster sube o reemplaza el póster y devuelve la URL pública.
func UploadPoster(key string, file multipart.File, mime string) (string, error) {
	if uploader == nil {
		return "", fmt.Errorf("el servicio de almacenamiento no está configurado")
	}

	fmt.Println("📤 [DEBUG] Subiendo a bucket:", bucket)
	fmt.Println("📤 [DEBUG] Key del archivo:", key)
	fmt.Println("📤 [DEBUG] Tipo MIME:", mime)

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configuración de carga específica para cada proveedor
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(mime),
	}

	// Solo añadir ACL si no es Supabase (que maneja permisos de forma diferente)
	if !isSupabase {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	// Intentar subir el archivo
	resp, err := uploader.Upload(ctx, input)
	if err != nil {
		fmt.Println("❌ [ERROR] Falló PutObject:", err)
		return "", err
	}

	fmt.Println("✅ [DEBUG] Subida exitosa. Location:", resp.Location)

	// Construye la URL pública
	url := fmt.Sprintf("%s/%s", baseURL, key)
	fmt.Println("✅ [DEBUG] URL pública:", url)
	return url, nil
}

// maskSecret oculta la mayor parte del secret para no imprimirlo completo.
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "********"
	}
	return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
}
