# Sistema de Puntuación de Comentarios (1-10)

## Descripción General

CineConecta implementa un sistema de análisis de sentimientos que evalúa el texto de los comentarios de los usuarios y asigna automáticamente una puntuación en escala de 1 a 10, donde:

- **1-3**: Sentimiento negativo (crítica)
- **4-6**: Sentimiento neutral
- **7-10**: Sentimiento positivo (recomendación)

## Funcionamiento del Análisis

El sistema analiza cada comentario de la siguiente manera:

1. **Normalización del texto**: Convierte todo a minúsculas y elimina caracteres especiales
2. **Análisis de palabras clave**: Identifica palabras positivas y negativas en el texto
3. **Detección de énfasis**: Reconoce modificadores de intensidad como "muy", "extremadamente", etc.
4. **Cálculo de puntuación**: Asigna un valor basado en la proporción de palabras positivas vs. negativas
5. **Conversión a escala 1-10**: Transforma el resultado a una escala intuitiva

## Ejemplos de Puntuación

| Puntuación | Descripción | Ejemplo de comentario |
|------------|-------------|----------------------|
| 9.5 - 10   | Obra maestra | "Una película absolutamente brillante, imperdible obra maestra del cine." |
| 8.0 - 9.4  | Excelente   | "Muy buena historia con actuaciones increíbles. Realmente la disfruté." |
| 7.0 - 7.9  | Muy buena   | "Me gustó mucho, bien dirigida y con buenos efectos especiales." |
| 6.0 - 6.9  | Buena       | "Entretenida aunque con algunos momentos lentos. Vale la pena verla." |
| 5.0 - 5.9  | Aceptable   | "Tiene sus momentos pero no es nada especial. Actuaciones decentes." |
| 4.0 - 4.9  | Regular     | "No está mal pero esperaba más. Historia predecible." |
| 3.0 - 3.9  | Mala        | "Bastante aburrida y con un guion pobre. No la recomendaría." |
| 2.0 - 2.9  | Muy mala    | "Una película terrible con actuaciones horribles. Pérdida de tiempo." |
| 1.0 - 1.9  | Pésima      | "Absolutamente horrible. Sin duda una de las peores películas que he visto." |

## Uso en Recomendaciones

Este sistema de puntuación forma la base de nuestro motor de recomendaciones:

1. **Recomendaciones personalizadas**: Identifica géneros preferidos basados en comentarios con puntuación ≥ 7.0
2. **Películas populares**: Clasifica películas según su puntuación promedio
3. **Filtrado por calidad**: Prioriza películas con buenas puntuaciones (≥ 6.0) en las recomendaciones

## API y Respuestas

Cuando se accede a la información de sentimiento de una película, la API devuelve:

```json
{
  "movie_id": 12,
  "sentiment": "positive",
  "sentiment_text": "positivo",
  "rating": 8.5,
  "rating_text": "Excelente"
}
```

Este enfoque proporciona tanto datos numéricos precisos como interpretaciones textuales fáciles de entender para los usuarios. 