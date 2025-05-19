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
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4" // ← firmante Sig-V4
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
	_ = godotenv.Load()

	endpointURL := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	region := strings.TrimSpace(os.Getenv("S3_REGION"))
	bucketEnv := strings.TrimSpace(os.Getenv("S3_BUCKET"))
	accessEnv := strings.TrimSpace(os.Getenv("S3_ACCESS_KEY"))
	secretEnv := strings.TrimSpace(os.Getenv("S3_SECRET_KEY"))

	fmt.Println("📦 [DEBUG] Configuración S3:")
	fmt.Println("    ENDPOINT :", endpointURL)
	fmt.Println("    REGION   :", region)
	fmt.Println("    BUCKET   :", bucketEnv)

	if endpointURL == "" || region == "" || bucketEnv == "" ||
		accessEnv == "" || secretEnv == "" {
		log.Println("⚠️  Variables S3 incompletas; Storage deshabilitado")
		return
	}

	bucket = bucketEnv
	isSupabase = strings.Contains(endpointURL, "supabase")

	if isSupabase {
		// URL pública para Supabase
		baseEndpoint := strings.Replace(endpointURL, "/storage/v1/s3", "", 1)
		baseURL = fmt.Sprintf("%s/storage/v1/object/public/%s", baseEndpoint, bucket)
	} else {
		baseURL = fmt.Sprintf("%s/%s", endpointURL, bucket)
	}

	// ---------- cliente S3 ----------
	var s3Client *s3.Client

	if isSupabase {
		resolver := s3.EndpointResolverFunc(
			func(region string, _ s3.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpointURL,
					PartitionID:   "aws",
					SigningRegion: region,
					SigningName:   "s3",
				}, nil
			})

		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessEnv, secretEnv, "")),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ cfg S3: %v", err))
		}

		s3Client = s3.NewFromConfig(
			cfg,
			s3.WithEndpointResolver(resolver),
			func(o *s3.Options) {
				o.UsePathStyle = true
				// ⚠️  Fuerza UNSIGNED-PAYLOAD (firma válida en Supabase)
				o.APIOptions = append(o.APIOptions,
					v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware)
			})
	} else {
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessEnv, secretEnv, "")),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ cfg S3: %v", err))
		}
		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
	}

	// ---------- uploader ----------
	uploader = manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.PartSize = 5 * 1024 * 1024 // 5 MB
		u.LeavePartsOnError = false  // 👈 evita la cabecera Content-MD5
	})
}

// UploadPoster sube el archivo y devuelve la URL pública.
func UploadPoster(key string, file multipart.File, mime string) (string, error) {
	if uploader == nil {
		return "", fmt.Errorf("almacenamiento no configurado")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(mime),
	}
	if !isSupabase {
		input.ACL = types.ObjectCannedACLPublicRead
	}

	if _, err := uploader.Upload(ctx, input); err != nil {
		return "", fmt.Errorf("error al subir póster: %w", err)
	}

	return fmt.Sprintf("%s/%s", baseURL, key), nil
}

// maskSecret ofusca la clave al imprimir
func maskSecret(s string) string {
	if len(s) <= 8 {
		return "********"
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}
