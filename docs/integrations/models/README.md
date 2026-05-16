# Plan: Módulo de Descarga de Modelos IA en Segundo Plano

TL;DR: Crear un módulo src/module/models/ con arquitectura hexagonal que descargue modelos de Hugging Face en paralelo (goroutine) sin bloquear el startup. Los metadatos se definen en data/models.local.go (Go struct), las descargas se guardan en data/models/, y se detectan descargas incompletas para reintentarlas. El módulo será expuesto como servicio para que assistant, voice, image puedan acceder a los modelos.

# Estructura de Fases

## Fase 1: Domain — Definir Contratos

src/module/models/domain/models.go — Structs:

Model — nombre, URL, versión, tamaño, ruta local
DownloadStatus — enum (Pending, Downloading, Complete, Failed)
ModelMetadata — lista de modelos disponibles
src/module/models/domain/usecase.go — Interfaz pura:

ModelLoaderUseCase interface con métodos: LoadModels(), ValidateDownload(), GetModelPath()

## Fase 2: Config — Metadatos de Modelos

data/models.local.go — Go struct con:
Array de modelos con nombre, URL Hugging Face, tamaño, checksum
Ejemplo:

## Fase 3: Infrastructure — Implementar IO

src/module/models/infrastructure/downloader.go — HuggingFaceDownloader:

Descargar archivo de HF con streaming
Guardar en data/models/{modelName}/
Crear .incomplete si se interrumpe
src/module/models/infrastructure/storage.go — LocalStorage:

CheckIfComplete() — verificar archivo .complete o tamaño
MarkComplete() — crear archivo .complete
CreateDirectoryIfNotExists() — crear data/models/

## Fase 4: Application — Lógica de Orquestación

src/module/models/application/loader.go — ModelLoaderUseCase:

LoadAllModels(ctx context.Context) — launch goroutines sin bloquear
ValidateAndRetryIncomplete() — detectar descargas fallidas
Loguear progreso/errores
src/module/models/application/service.go — ModelService:

Exponer métodos públicos para otros módulos: GetModelPath(name), IsModelReady(name)

## Fase 5: Integración en main.go

En main.go alrededor de // 4.1 CARGAR MODELOS DE IA LOCALES:
Crear instancia de ModelLoaderUseCase
Llamar go loader.LoadAllModels(ctx) para no bloquear
Guardar ModelService en contexto global o singleton para otros módulos

## Fase 6: Testing

test/unit/module/models/domain_test.go — Tests domain:

Validar structs, enums
test/unit/module/models/downloader_test.go — Tests downloader:

Mock Hugging Face responses
Simular descarga completa e incompleta
test/unit/module/models/storage_test.go — Tests storage:

Crear/validar archivos .complete, .incomplete
test/e2e/module/models/loader_test.go — Tests e2e:

Descarga real (o mock) y validación de carpetas
Archivos a Crear/Modificar
Nuevos:

src/module/models/domain/models.go
src/module/models/domain/usecase.go
src/module/models/application/loader.go
src/module/models/application/service.go
src/module/models/infrastructure/downloader.go
src/module/models/infrastructure/storage.go
src/module/models/module.go (actualmente vacío)
data/models.local.go
Tests en test/unit/module/models/ y test/e2e/module/models/
Modificar:

main.go — Agregar inicialización de ModelLoader en // 4.1 CARGAR MODELOS DE IA LOCALES

Verificación
✅ go mod tidy && go mod download sin errores
✅ Carpeta data/models/ existe y es creada automáticamente
✅ Al arrancar, logs muestran "Descargando modelos..." sin bloquear
✅ Tests pasan: go test ./src/module/models/...
✅ Archivo .complete en data/models/{modelName}/ cuando termine
✅ Reintentar si encuentra .incomplete al arrancar
Decisiones
✅ Formato: Go struct (data/models.local.go) para eficiencia y type-safety
✅ Goroutines: go loader.LoadAllModels() sin bloquear main
✅ Errores: Log + continue, con detección de incompletos para reintentos
✅ Localización: src/module/models/ (no en src/http/ porque no expone rutas HTTP aún)
✅ Service expuesto: Otros módulos accederán via ModelService.GetModelPath()
