# LLM Local en Kara (Implementacion Actual)

Este documento resume como funciona la inferencia local con llama.cpp en apps/core, como se integra con WebSocket y que debes considerar para operar esto en desarrollo y en produccion.

## 1. Alcance de la implementacion

Actualmente el backend de Go en apps/core implementa:

- Descarga en segundo plano de modelos GGUF desde Hugging Face.
- Descarga y preparacion en segundo plano de binarios llama.cpp por plataforma.
- Ejecucion de inferencia local usando llama-server (modo persistente).
- Integracion por WebSocket en los eventos message y chat.
- Persistencia basica de mensajes user/assistant en SQLite.
- Apagado limpio del proceso HTTP y del proceso llama-server.

## 2. Arquitectura

Capas principales:

- Modulo models:
    - Domain/Application/Infrastructure para gestionar modelos GGUF locales.
    - Metadatos en data/models.local.go.
- Modulo inference:
    - Downloader de binarios en data/binaries.local.go.
    - Executor local que levanta llama-server y llama /v1/chat/completions.
    - Service con estado ready/not-ready.
- Modulo websocket:
    - Evento message: reply directo del modelo.
    - Evento chat: respuesta markdown chunked.
    - Guardado de intercambio en tabla messages.

Orquestacion:

- main.go inicia SQLite.
- Lanza descargas de modelos en goroutine.
- Lanza preparacion de binario de inferencia en goroutine.
- Registra rutas HTTP y WS.
- En shutdown cierra Fiber y luego cierra inference service.

## 3. Flujo de arranque

1. Se carga configuracion (env + defaults).
2. Se abre base de datos, migraciones, validacion, seed.
3. Se dispara carga de modelos locales (no bloqueante).
4. Se dispara preparacion de llama.cpp (no bloqueante).
5. Cuando binario esta listo, inference service pasa a ready=true.
6. WebSocket puede generar respuestas si service IsReady().

Nota importante:

- Si inference no esta ready, los handlers retornan error de binario no listo.

## 4. Configuracion de entorno (LLM)

Variables relevantes en apps/core/src/server/config/config.go:

- INFERENCE_ENABLED
    - Default: true
    - Si false, no se prepara inference.
- INFERENCE_MODEL
    - Default: qwen3-1.7b-light
    - Modelo usado cuando el request no define uno.
- INFERENCE_TIMEOUT_SECONDS
    - Default: 120
    - Timeout HTTP para requests al llama-server.
- INFERENCE_MAX_TOKENS
    - Default: 256
    - Max tokens por request si no se especifica en runtime.
- INFERENCE_TEMPERATURE
    - Default: 0.7
    - Temperatura por defecto.
- INFERENCE_SERVER_HOST
    - Default: 127.0.0.1
    - Host local para bind de llama-server.
- INFERENCE_SERVER_PORT
    - Default: 32111
    - Puerto preferido para iniciar llama-server.
- INFERENCE_SERVER_PORT_MIN
    - Default: 32111
    - Inicio del rango de fallback de puertos.
- INFERENCE_SERVER_PORT_MAX
    - Default: 32130
    - Fin del rango de fallback de puertos.

Recomendacion inicial para estabilidad local:

- INFERENCE_TIMEOUT_SECONDS=120
- INFERENCE_MAX_TOKENS=256 o 384
- INFERENCE_TEMPERATURE=0.3 a 0.7
- INFERENCE_SERVER_HOST=127.0.0.1
- INFERENCE_SERVER_PORT=32111
- INFERENCE_SERVER_PORT_MIN=32111
- INFERENCE_SERVER_PORT_MAX=32130

## 5. Modelos y binarios definidos

Modelos GGUF (data/models.local.go):

- qwen3-1.7b-light
- qwen3-4b-instruct
- qwen2.5-coder-7b-instruct

Binarios llama.cpp (data/binaries.local.go):

- Version b9174
- Ejecutable objetivo: llama-server(.exe)
- Assets por os/arch para windows/linux/darwin

Politica de puertos anti-colision:

- Puerto preferido configurable con INFERENCE_SERVER_PORT.
- Si el puerto preferido esta ocupado, se intenta fallback automatico en el rango definido por INFERENCE_SERVER_PORT_MIN y INFERENCE_SERVER_PORT_MAX.
- Bind solo en loopback usando INFERENCE_SERVER_HOST (por defecto 127.0.0.1).
- El arranque registra en logs el puerto final seleccionado y si se uso fallback.

## 6. Integracion WebSocket

Endpoint:

- GET /ws

Eventos entrantes:

- message
    - data puede ser string o objeto con message.
- chat
    - data con forma { message, context? }.

Eventos salientes relevantes:

- connected
- message (echo)
- message:reply
- chat:chunk
- chat:done
- chat:accepted
- error

Comportamiento:

- message devuelve content del modelo en una sola respuesta.
- chat devuelve markdown en chunks y cierra con done/accepted.

## 7. Manejo especial de respuestas vacias en Qwen

Problema observado:

- llama-server puede devolver choices[0].message.content vacio.
- A la vez puede traer reasoning_content lleno y finish_reason=length.

Mitigacion implementada:

- Se detecta el patron content vacio + finish_reason=length + reasoning presente.
- Se hace un reintento automatico una sola vez con:
    - prompt reforzado a respuesta final sin razonamiento interno.
    - mas tokens (con tope).
    - temperatura mas baja (maximo 0.3).

Objetivo:

- Reducir fallos en chat por consumo de tokens en cadena de razonamiento.

## 8. Persistencia en base de datos

Al generar respuesta valida:

- Inserta mensaje user en messages.
- Inserta mensaje assistant en messages con latency_ms.
- Si no hay user/assistant/session, se crean por defecto.

Esto permite mantener historial basico del intercambio WS.

## 9. Apagado limpio

En cierre del proceso:

- Fiber recibe shutdown con timeout.
- Luego se ejecuta inferenceModule.Service.Shutdown().
- LocalExecutor cierra llama-server con senal y kill de respaldo.

Beneficio:

- Menos procesos huerfanos y menos riesgo de puerto ocupado al reiniciar.

## 10. Como ejecutarlo

Desde root:

- npm run start
- npm run start:core

Desde apps/core:

- npm run start
- go run main.go

Validacion tecnica:

- go test ./...

## 11. Checklist operativo

Antes de probar chat local:

1. Confirmar que SQLite arranco sin errores.
2. Confirmar logs de preparacion de llama.cpp.
3. Confirmar que inference service quede en ready.
4. Enviar evento message por WS y validar message:reply.
5. Enviar evento chat y validar chunk/done/accepted.

## 12. Riesgos y limites actuales

- El rango de fallback puede agotarse en entornos muy cargados.
- Un solo proceso de inferencia local por instancia.
- No hay cola avanzada ni scheduler de prompts.
- No hay auth en /ws (CheckOrigin abierto en true).
- El fallback de chat sigue generando markdown cuando inference falla.

## 13. Mejoras recomendadas

1. Hacer puerto de llama-server configurable por env.
2. Permitir politicas de exclusion de puertos por entorno (lista denylist configurable).
3. Agregar metrica explicita de retries por reasoning.
4. Endurecer seguridad de WS (origin, auth, limites).
5. Parametrizar estrategia de retry (on/off, max tokens, temp).

---

Documentos relacionados:

- docs/integrations/models/README.md
- docs/llm/windows.md
- docs/llm/troubleshooting.md
