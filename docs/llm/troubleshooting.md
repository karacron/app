# Troubleshooting LLM Local

Guia de fallas comunes y respuestas operativas para la implementacion actual.

## 1. Error: binary not ready

Sintoma:

- message/chat responde con error de inferencia no disponible.

Causa probable:

- Downloader de llama.cpp fallo o sigue en proceso.

Que hacer:

1. Revisar logs de arranque.
2. Confirmar descarga correcta en data/binaries.
3. Reiniciar servicio despues de validar archivos.

## 2. Error: llama-server no disponible

Sintoma:

- Fallo HTTP al llamar /v1/chat/completions.

Causa probable:

- Proceso no iniciado, fallo en startup o puerto bloqueado.

Que hacer:

1. Revisar logs de ensureServer.
2. Revisar INFERENCE_SERVER_PORT y confirmar que el host es loopback (INFERENCE_SERVER_HOST).
3. Revisar INFERENCE_SERVER_PORT_MIN e INFERENCE_SERVER_PORT_MAX.
4. Verificar si se agotaron puertos del rango de fallback configurado.
5. Confirmar ejecutable y permisos de ejecucion.

## 3. Error: respondio sin contenido

Sintoma:

- body con choices pero sin message.content.

Caso especial Qwen:

- message.reasoning_content viene con texto.
- finish_reason llega como length.

Mitigacion ya implementada:

- Retry automatico con prompt final-only, mas tokens y menor temperatura.

Si persiste:

1. Aumentar INFERENCE_MAX_TOKENS.
2. Usar prompts mas directos.
3. Probar con temperatura mas baja.
4. Cambiar a un modelo menos propenso a razonamiento largo para esa tarea.

## 4. Error 0xc0000135 en Windows

Sintoma:

- El proceso termina al iniciar.

Causa:

- DLL runtime faltantes.

Que hacer:

1. Revalidar descarga y extraccion del asset.
2. Verificar que runtime completo este junto al ejecutable.
3. Reiniciar backend.

## 5. Error no se puede encontrar la ruta especificada

Causa:

- Uso de rutas relativas invalidas en ejecucion.

Que hacer:

1. Usar rutas absolutas para executable y modelPath.
2. Confirmar existencia fisica de archivos.

## 6. Chat responde fallback markdown en vez de modelo

Causa:

- inferPrompt devolvio error y se activo fallback.

Que hacer:

1. Revisar evento de error asociado.
2. Corregir causa raiz de inference.
3. Reprobar con evento message para validar camino corto.

## 7. No persiste mensajes en SQLite

Causa probable:

- Error al asegurar actores/sesion o al insertar.

Que hacer:

1. Revisar logs ws de persistExchange.
2. Validar migraciones y schema de tablas.
3. Confirmar DATABASE_PATH correcto y writable.

## 8. Pruebas recomendadas tras cambios

1. go test ./...
2. Smoke test WS con message.
3. Smoke test WS con chat.
4. Verificacion de inserciones en messages.
