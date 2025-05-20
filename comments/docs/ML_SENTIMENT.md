# Análisis de Sentimientos con Machine Learning

Este módulo implementa un análisis de sentimientos avanzado para los comentarios de películas utilizando la API de OpenAI, permitiendo un análisis más preciso y contextual de los comentarios en español.

## Configuración

El sistema soporta dos modos de análisis de sentimientos:

1. **Modo Básico (Léxico):** Análisis basado en diccionario de palabras positivas/negativas, que funciona sin dependencias externas.
2. **Modo ML:** Análisis avanzado utilizando el modelo GPT de OpenAI, que proporciona resultados más precisos entendiendo el contexto y la semántica.

### Requisitos para el modo ML

- API Key de OpenAI configurada
- Variable de entorno `USE_ML_SENTIMENT=true`

## API de Configuración

### Obtener configuración actual
```
GET /api/comments/settings (AdminRequired)
```

Respuesta:
```json
{
  "use_ml": true,
  "has_openai_key": true
}
```

### Actualizar configuración
```
POST /api/comments/settings (AdminRequired)
```

Cuerpo de la petición:
```json
{
  "use_ml": true,
  "openai_key": "sk-xxxxxxxxxx" 
}
```

## Funcionamiento

1. Cuando se crea o actualiza un comentario, el sistema determina si debe usar ML o el análisis léxico según la configuración.
2. Si se usa ML, se envía el comentario a la API de OpenAI con instrucciones específicas para clasificarlo.
3. La API devuelve un JSON con:
   - `sentiment`: Clasificación (positive, neutral, negative)
   - `score`: Puntuación en escala 1-10
   - `reason`: Explicación del análisis (no se almacena, solo informativo)
4. Si falla la API o no está configurada, se usa el método léxico como fallback.

## Re-procesar comentarios existentes

Puedes re-analizar todos los comentarios existentes con el nuevo método:

```
POST /api/comments/update-sentiments (AdminRequired)
```

## Consideraciones técnicas

- El modo ML consume créditos de la API de OpenAI.
- Para comentarios muy cortos (<10 caracteres), se usa siempre el análisis léxico.
- Se utiliza una temperatura de 0.0 para resultados consistentes y deterministas.
- La respuesta está limitada a 150 tokens para minimizar costos. 