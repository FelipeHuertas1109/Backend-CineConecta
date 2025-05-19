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
	endpointURL := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	region := strings.TrimSpace(os.Getenv("S3_REGION"))
	bucketEnv := strings.TrimSpace(os.Getenv("S3_BUCKET"))
	accessEnv := strings.TrimSpace(os.Getenv("S3_ACCESS_KEY"))
	secretEnv := strings.TrimSpace(os.Getenv("S3_SECRET_KEY"))

	fmt.Println("📦 [DEBUG] Configuración S3 cargada:")
	fmt.Println("    S3_ENDPOINT  :", endpointURL)
	fmt.Println("    S3_REGION    :", region)
	fmt.Println("    S3_BUCKET    :", bucketEnv)
	fmt.Println("    S3_ACCESS_KEY:", maskSecret(accessEnv))
	fmt.Println("    S3_SECRET_KEY:", maskSecret(secretEnv))

	if endpointURL == "" || region == "" || bucketEnv == "" || accessEnv == "" || secretEnv == "" {
		log.Println("⚠️ [ADVERTENCIA] Variables de entorno S3 incompletas. Subida de imágenes deshabilitada.")
		return
	}

	// 3) Asigna variables globales
	bucket = bucketEnv
	isSupabase = strings.Contains(endpointURL, "supabase")

	// Para Supabase, la URL pública debe construirse correctamente
	if isSupabase {
		// https://<proyecto>.supabase.co/storage/v1/object/public/<bucket>/<path>
		baseEndpoint := strings.Replace(endpointURL, "/storage/v1/s3", "", 1)
		baseURL = fmt.Sprintf("%s/storage/v1/object/public/%s", baseEndpoint, bucket)
	} else {
		baseURL = fmt.Sprintf("%s/%s", endpointURL, bucket)
	}

	fmt.Println("    BASE_URL     :", baseURL)
	fmt.Println("    ES SUPABASE  :", isSupabase)

	// 4) Configuración del cliente S3
	var s3Client *s3.Client

	if isSupabase {
		// Resolver específico para Supabase (path-style + servicio "s3")
		supabaseResolver := s3.EndpointResolverFunc(
			func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpointURL, // p.ej. https://zufjxpgxyhphoygtxqit.supabase.co/storage/v1/s3
					PartitionID:   "aws",       // requerido por el SDK
					SigningRegion: region,
					SigningName:   "s3", // ← clave para que la firma coincida
				}, nil
			})

		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessEnv, secretEnv, ""),
			),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Error al cargar configuración S3: %v", err))
		}

		// Cliente S3 para Supabase: path-style y resolver propio
		s3Client = s3.NewFromConfig(cfg,
			s3.WithEndpointResolver(supabaseResolver),
			func(o *s3.Options) { o.UsePathStyle = true },
		)
	} else {
		// Cliente estándar para otros proveedores
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessEnv, secretEnv, "")),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Error al cargar configuración S3: %v", err))
		}

		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	// 5) Inicializa el uploader
	uploader = manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.PartSize = 5 * 1024 * 1024 // 5 MB
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(mime),
	}

	// Supabase maneja permisos de manera interna
	if !isSupabase {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	resp, err := uploader.Upload(ctx, input)
	if err != nil {
		fmt.Println("❌ [ERROR] Falló PutObject:", err)
		return "", err
	}

	fmt.Println("✅ [DEBUG] Subida exitosa. Location:", resp.Location)

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
