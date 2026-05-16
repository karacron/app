# Operacion en Windows (LLM local)

Este documento concentra consideraciones especificas de Windows para la inferencia local con llama.cpp.

## 1. Error 0xc0000135 (DLL faltantes)

Sintoma:

- El ejecutable no inicia y retorna 0xc0000135.

Causa habitual:

- Solo se copia llama-server.exe sin las DLL runtime del release.

Accion requerida:

- Copiar runtime completo del asset extraido al directorio final de ejecucion.
- No asumir que con el .exe alcanza.

## 2. Rutas de ejecucion

Buenas practicas ya aplicadas:

- Resolver ruta absoluta del ejecutable y del modelo GGUF.
- Ejecutar con cmd.Dir en el directorio del binario.

Por que:

- Evita fallos tipo no se puede encontrar la ruta especificada.

## 3. Señales y apagado

Comportamiento esperado:

- En shutdown se envia os.Interrupt.
- Si no cierra en tiempo razonable, se fuerza kill.

Observacion:

- En Windows puedes ver codigos de salida por interrupcion de consola; no siempre implican corrupcion.

## 4. Antivirus y SmartScreen

Riesgos operativos:

- Algunos antivirus pueden bloquear binarios descargados o eliminar DLL.
- SmartScreen puede advertir en primera ejecucion.

Recomendaciones:

1. Mantener carpeta de data/binaries en exclusions si el entorno lo requiere.
2. Verificar integridad de archivos tras descarga/extraccion.
3. Registrar logs detallados cuando la preparacion falle.

## 5. Rendimiento y consumo

Consejos para equipos Windows de escritorio:

- Para baja RAM, usar modelo mas pequeno (qwen3-1.7b-light).
- Reducir INFERENCE_MAX_TOKENS si hay latencias altas.
- Evitar multiples procesos que compitan por CPU al mismo tiempo.

## 6. Red local y puerto

Estado actual:

- llama-server corre en 127.0.0.1:8080.

Puntos de cuidado:

- Si otro proceso usa 8080, inference puede no iniciar.
- Debe evitarse exponer ese puerto fuera de localhost.

## 7. Diagnostico rapido

Si no responde el modelo:

1. Revisar logs de preparacion del binario en arranque.
2. Confirmar que la ruta GGUF exista en data/models.
3. Confirmar que service este en ready.
4. Revisar si hay error de puerto ocupado.
5. Probar de nuevo con prompt corto en evento message.
