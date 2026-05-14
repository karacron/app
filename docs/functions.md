# Electron Background Assistant Mode

## Objetivo

Implementar en la aplicación Electron un modo de ejecución en segundo plano tipo asistente local, donde la app pueda permanecer activa aunque el usuario cierre la ventana principal. En Windows, la aplicación debe quedar disponible en el System Tray, junto al reloj del sistema, y permitir mostrar una mini ventana flotante mediante una combinación global de teclas.

## Requisitos funcionales

### 1. Comportamiento al cerrar la aplicación

Cuando el usuario pulse cerrar en la ventana principal:

- No debe finalizar el proceso de Electron.
- La ventana principal debe ocultarse usando `window.hide()`.
- La aplicación debe seguir ejecutándose en segundo plano.
- El usuario debe poder restaurarla desde el icono del System Tray.
- Debe existir una opción explícita de "Salir" en el menú del tray para cerrar realmente la aplicación.

### 2. System Tray en Windows

Implementar un icono persistente en el área de notificación de Windows usando `Tray`.

El menú contextual del tray debe incluir:

- Abrir asistente
- Mostrar mini chat
- Ocultar ventana
- Estado del asistente
- Iniciar con Windows
- Salir

El icono debe tener tooltip con el nombre de la aplicación.

### 3. Ejecución en segundo plano

La lógica crítica del asistente no debe depender de que una ventana esté visible.

El `main process` debe gestionar:

- Estado global de la aplicación.
- Ciclo de vida del tray.
- Registro de atajos globales.
- Comunicación IPC con ventanas.
- Workers o procesos secundarios si hay tareas pesadas.
- Persistencia de estado.
- Eventos de bloqueo/desbloqueo del sistema.

La UI React/renderer solo debe encargarse de mostrar la interfaz y enviar acciones al proceso principal mediante IPC seguro.

### 4. Atajo global

Registrar una combinación global de teclas con `globalShortcut`.

Ejemplo recomendado:

- `CommandOrControl+Shift+Space`

Cuando el usuario pulse el atajo:

- Si la mini ventana no existe, crearla.
- Si existe y está visible, ocultarla.
- Si existe y está oculta, mostrarla.
- La ventana debe aparecer centrada o cerca de la parte inferior derecha.
- Debe ser pequeña, rápida y enfocada al chat con el asistente.

### 5. Mini ventana flotante del asistente

Crear una `BrowserWindow` secundaria para el modo asistente rápido.

Características:

- Tamaño aproximado: 420x560.
- Sin frame nativo o con diseño minimal.
- Siempre encima cuando esté visible usando `alwaysOnTop`.
- Ocultable al perder foco, si aplica.
- No debe aparecer en la taskbar si se comporta como overlay.
- Debe cargar una ruta específica de la app, por ejemplo `/assistant-mini`.
- Debe comunicarse con el backend local o el main process mediante IPC seguro.

### 6. Bloqueo y desbloqueo de Windows

Usar `powerMonitor` para detectar:

- `lock-screen`
- `unlock-screen`
- `suspend`
- `resume`

Comportamiento esperado:

- Al bloquear pantalla:
    - Ocultar ventanas visibles.
    - Mantener tareas internas si el sistema sigue activo.
    - Pausar únicamente tareas que requieran interacción visual.
- Al desbloquear:
    - Restaurar estado interno.
    - Revalidar conexiones locales.
    - Reanudar tareas pausadas.
    - No mostrar ventanas automáticamente salvo que el usuario lo haya configurado.

Importante: Windows no permite mostrar una ventana Electron por encima de la pantalla de bloqueo. La app puede seguir ejecutando lógica en segundo plano dentro de la sesión del usuario mientras el equipo no esté suspendido, pero la UI solo debe mostrarse cuando la sesión esté desbloqueada.

### 7. Inicio automático con Windows

Permitir que el usuario active/desactive el inicio automático con Windows.

Usar:

```ts
app.setLoginItemSettings({
	openAtLogin: true,
});

type AppSettings = {
	openAtLogin: boolean;
	startMinimizedToTray: boolean;
	enableGlobalShortcut: boolean;
	globalShortcut: string;
};
```

## Si la app arranca con Windows y startMinimizedToTray está activo:

No mostrar la ventana principal.
Inicializar tray.
Inicializar servicios en segundo plano.
Registrar atajos globales.

# 8. Seguridad

Aplicar buenas prácticas de seguridad en Electron:

contextIsolation: true
nodeIntegration: false
sandbox: true cuando sea posible.
Usar preload.ts para exponer una API mínima.
Validar todos los mensajes IPC.
No exponer APIs internas de Node.js al renderer.
No ejecutar comandos arbitrarios desde el renderer.

## Separar claramente:

- main
- preload
- renderer
- background services

# 10. Ejemplo técnico base

Implementar un flujo parecido a este:

```ts
import {
	app,
	BrowserWindow,
	Tray,
	Menu,
	globalShortcut,
	powerMonitor,
} from "electron";
import path from "node:path";

let mainWindow: BrowserWindow | null = null;
let assistantWindow: BrowserWindow | null = null;
let tray: Tray | null = null;
let isQuitting = false;

function createMainWindow() {
	mainWindow = new BrowserWindow({
		width: 1200,
		height: 800,
		show: false,
		webPreferences: {
			preload: path.join(__dirname, "../preload/index.js"),
			contextIsolation: true,
			nodeIntegration: false,
		},
	});

	mainWindow.loadURL("app://index");

	mainWindow.once("ready-to-show", () => {
		mainWindow?.show();
	});

	mainWindow.on("close", (event) => {
		if (!isQuitting) {
			event.preventDefault();
			mainWindow?.hide();
		}
	});
}

function createAssistantWindow() {
	assistantWindow = new BrowserWindow({
		width: 420,
		height: 560,
		show: false,
		frame: false,
		resizable: false,
		alwaysOnTop: true,
		skipTaskbar: true,
		webPreferences: {
			preload: path.join(__dirname, "../preload/index.js"),
			contextIsolation: true,
			nodeIntegration: false,
		},
	});

	assistantWindow.loadURL("app://assistant-mini");

	assistantWindow.on("blur", () => {
		assistantWindow?.hide();
	});
}

function toggleAssistantWindow() {
	if (!assistantWindow) {
		createAssistantWindow();
	}

	if (assistantWindow?.isVisible()) {
		assistantWindow.hide();
		return;
	}

	assistantWindow?.show();
	assistantWindow?.focus();
}

function createTray() {
	tray = new Tray(path.join(__dirname, "../../assets/tray.ico"));

	const menu = Menu.buildFromTemplate([
		{
			label: "Abrir aplicación",
			click: () => {
				mainWindow?.show();
				mainWindow?.focus();
			},
		},
		{
			label: "Mostrar asistente",
			click: () => toggleAssistantWindow(),
		},
		{
			type: "separator",
		},
		{
			label: "Salir",
			click: () => {
				isQuitting = true;
				app.quit();
			},
		},
	]);

	tray.setToolTip("Local AI Assistant");
	tray.setContextMenu(menu);

	tray.on("double-click", () => {
		mainWindow?.show();
		mainWindow?.focus();
	});
}

function registerGlobalShortcuts() {
	globalShortcut.register("CommandOrControl+Shift+Space", () => {
		toggleAssistantWindow();
	});
}

function registerPowerEvents() {
	powerMonitor.on("lock-screen", () => {
		mainWindow?.hide();
		assistantWindow?.hide();
	});

	powerMonitor.on("unlock-screen", () => {
		// Reanudar servicios internos si aplica.
	});

	powerMonitor.on("suspend", () => {
		// Pausar tareas no críticas.
	});

	powerMonitor.on("resume", () => {
		// Revalidar conexiones y reanudar servicios.
	});
}

app.whenReady().then(() => {
	app.setLoginItemSettings({
		openAtLogin: true,
	});

	createMainWindow();
	createTray();
	registerGlobalShortcuts();
	registerPowerEvents();
});

app.on("window-all-closed", (event) => {
	event.preventDefault();
});

app.on("before-quit", () => {
	isQuitting = true;
});

app.on("will-quit", () => {
	globalShortcut.unregisterAll();
});
```

# 11. Modo servicio real por sistema operativo

Para garantizar que la aplicación siempre esté viva, separar responsabilidades en dos procesos:

- Servicio de backend (siempre vivo): `kara-core` (HTTP/IPC local).
- UI Electron (opcional): se puede abrir/cerrar sin afectar al servicio.

La UI debe reconectar automáticamente al backend al mostrarse.

## 11.1 Arquitectura recomendada

1. `kara-core` corre como servicio del sistema.
2. Electron se conecta a `http://127.0.0.1:8081` (o socket local) para estado/comandos.
3. Si el usuario cierra la ventana, solo se oculta la UI.
4. Si Electron se cierra por completo, el backend sigue activo.
5. Al volver a abrir Electron, se reanuda la sesión contra el backend.

## 11.2 macOS (launchd)

### Opción A: servicio por usuario (recomendado para desktop)

Archivo: `~/Library/LaunchAgents/com.kara.core.plist`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>Label</key><string>com.kara.core</string>
		<key>ProgramArguments</key>
		<array>
			<string>/opt/kara/kara-core</string>
		</array>
		<key>WorkingDirectory</key><string>/opt/kara</string>
		<key>RunAtLoad</key><true/>
		<key>KeepAlive</key><true/>
		<key>StandardOutPath</key><string>/tmp/kara-core.out.log</string>
		<key>StandardErrorPath</key><string>/tmp/kara-core.err.log</string>
	</dict>
</plist>
```

Comandos:

```bash
launchctl unload ~/Library/LaunchAgents/com.kara.core.plist 2>/dev/null || true
launchctl load ~/Library/LaunchAgents/com.kara.core.plist
launchctl start com.kara.core
```

### Opción B: servicio global

Usar `/Library/LaunchDaemons` cuando necesites arrancar sin sesión de usuario. En ese caso, la UI Electron debe iniciarse aparte al login del usuario.

## 11.3 Linux (systemd)

Archivo: `/etc/systemd/system/kara-core.service`

```ini
[Unit]
Description=Kara Core Service
After=network.target

[Service]
Type=simple
User=kara
WorkingDirectory=/opt/kara
ExecStart=/opt/kara/kara-core
Restart=always
RestartSec=3
Environment=KARA_PORT=8081

[Install]
WantedBy=multi-user.target
```

Comandos:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now kara-core
sudo systemctl status kara-core
```

Logs:

```bash
journalctl -u kara-core -f
```

## 11.4 Windows (Windows Service)

Punto clave: un servicio Windows corre en Session 0 y no puede mostrar UI (tray/ventanas). Por diseño, UI y servicio deben ir separados.

### Servicio del backend

Registrar `kara-core.exe` como servicio con una de estas opciones:

- `sc.exe` (nativo)
- NSSM (recomendado para reinicios y logs)
- WinSW (wrapper XML)

Ejemplo con `sc.exe`:

```powershell
sc.exe create KaraCore binPath= "C:\Program Files\Kara\kara-core.exe" start= auto
sc.exe failure KaraCore reset= 86400 actions= restart/5000/restart/5000/restart/5000
sc.exe start KaraCore
```

### Inicio de la UI Electron

Para tray, atajos globales y mini ventana, iniciar Electron en sesión de usuario con:

- `app.setLoginItemSettings({ openAtLogin: true })`
- o Task Scheduler (al iniciar sesión del usuario)

No ejecutar Electron principal como Windows Service si necesitas interacción de escritorio.

## 11.5 Contrato entre UI y servicio

Definir un contrato mínimo para desacoplar frontend y backend:

- `GET /health` para disponibilidad.
- `GET /state` para estado global.
- `POST /assistant/toggle` o IPC equivalente para acciones.

Comportamiento de la UI:

1. Si `/health` responde, operar normalmente.
2. Si falla, mostrar "Reconectando" y reintentar con backoff.
3. Nunca bloquear el proceso de UI esperando una respuesta síncrona del backend.

## 11.6 Inicio y recuperación

- Arranque del sistema: se levanta `kara-core`.
- Login de usuario: arranca Electron minimizado a tray (si está habilitado).
- Caída de backend: el supervisor lo reinicia (`KeepAlive`/`Restart=always`/recovery actions).
- Caída de UI: no afecta al backend; la UI puede abrirse después y reconectar.

## 11.7 Seguridad operativa del servicio

- Ejecutar con usuario de menor privilegio.
- Restringir bind a localhost (`127.0.0.1`) salvo necesidad explícita.
- Firmar binarios (especialmente en Windows/macOS).
- Rotar logs y evitar secretos en texto plano.
- Añadir watchdog de salud y métricas básicas.

## 11.8 Automatizacion obligatoria al instalar la app

Objetivo: cuando el usuario instala la app, el backend `kara-core` debe quedar registrado y arrancado como servicio sin pasos manuales.

### Flujo que debemos implementar nosotros

1. Instalacion:
    - Copiar binario `kara-core` a ruta estable.
    - Registrar servicio del SO (launchd/systemd/Windows Service).
    - Habilitar arranque automatico.
    - Iniciar servicio y validar `GET /health`.
2. Actualizacion:
    - Detener servicio.
    - Reemplazar binario.
    - Reaplicar definicion del servicio (idempotente).
    - Iniciar y validar salud.
3. Desinstalacion:
    - Detener y desregistrar servicio.
    - Eliminar archivos de servicio y binario.

Regla clave: todos los scripts deben ser idempotentes (si ya existe, actualizar; si no existe, crear; si no esta activo, iniciar).

## 11.9 Implementacion automatizada por plataforma

### 11.9.1 Windows

Recomendado en instalador:

- Empaquetar con electron-builder + NSIS.
- Ejecutar script elevado en `customInstall` para crear/actualizar servicio.
- Ejecutar script en `customUnInstall` para eliminar servicio.

Acciones de script:

```powershell
# install-service.ps1
$serviceName = "KaraCore"
$exePath = "C:\Program Files\Kara\kara-core.exe"

if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
	sc.exe stop $serviceName | Out-Null
	sc.exe config $serviceName binPath= "`"$exePath`"" start= auto | Out-Null
} else {
	sc.exe create $serviceName binPath= "`"$exePath`"" start= auto | Out-Null
}

sc.exe failure $serviceName reset= 86400 actions= restart/5000/restart/5000/restart/5000 | Out-Null
sc.exe start $serviceName | Out-Null
```

```powershell
# uninstall-service.ps1
$serviceName = "KaraCore"
if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
	sc.exe stop $serviceName | Out-Null
	sc.exe delete $serviceName | Out-Null
}
```

Nota: Electron (tray/mini chat) se inicia por login de usuario (`setLoginItemSettings`), no como servicio.

### 11.9.2 macOS

Recomendado en instalador:

- Para automatizacion real, distribuir `.pkg` con scripts `postinstall` y `preinstall`.
- `postinstall` copia/actualiza plist en `~/Library/LaunchAgents` (por usuario) o `/Library/LaunchDaemons` (global), luego `launchctl load`.

Acciones de script:

```bash
# postinstall (ejemplo LaunchAgent por usuario actual)
PLIST="$HOME/Library/LaunchAgents/com.kara.core.plist"
mkdir -p "$HOME/Library/LaunchAgents"
cp "/Applications/Kara.app/Contents/Resources/service/com.kara.core.plist" "$PLIST"

launchctl unload "$PLIST" 2>/dev/null || true
launchctl load "$PLIST"
launchctl start com.kara.core || true
```

```bash
# preinstall o uninstall helper
PLIST="$HOME/Library/LaunchAgents/com.kara.core.plist"
launchctl unload "$PLIST" 2>/dev/null || true
rm -f "$PLIST"
```

Nota: con `.dmg` puro no hay ciclo de postinstall robusto; para servicio automatizado usar `.pkg`.

### 11.9.3 Linux

Recomendado en instalador:

- Distribuir `.deb`/`.rpm` (no solo AppImage) para tener hooks del sistema.
- Usar `postinst`/`prerm`/`postrm` para registrar y retirar `systemd`.

Acciones de script (`.deb`):

```bash
# postinst
install -m 755 /opt/kara/kara-core /opt/kara/kara-core
install -m 644 /opt/kara/kara-core.service /etc/systemd/system/kara-core.service
systemctl daemon-reload
systemctl enable --now kara-core
```

```bash
# prerm/postrm
systemctl stop kara-core || true
systemctl disable kara-core || true
rm -f /etc/systemd/system/kara-core.service
systemctl daemon-reload
```

Nota: AppImage no instala servicio del sistema por defecto; si solo se distribuye AppImage, incluir un asistente de primera ejecucion que haga instalacion guiada con permisos.

## 11.10 Checklist de release (obligatorio)

Antes de publicar cada version:

1. Verificar instalacion limpia: servicio creado y activo.
2. Verificar actualizacion sobre version anterior: servicio preservado.
3. Verificar desinstalacion: servicio eliminado.
4. Verificar rollback (si aplica): servicio vuelve a version previa.
5. Verificar que Electron abre/cierra sin tumbar `kara-core`.
6. Verificar logs y estado: `health` responde tras reinicio de equipo.

## 11.11 Proceso de actualizacion automatizada

Si, este proceso se puede y se recomienda validarlo con GitHub Releases.

### Flujo recomendado (end-to-end)

1. CI compila artefactos por plataforma (`.exe/.msi`, `.pkg/.dmg`, `.deb/.rpm`).
2. CI publica release en GitHub con version semantica (`vX.Y.Z`).
3. La app Electron consulta proveedor de actualizaciones (GitHub provider).
4. Si hay nueva version:
    - descarga paquete,
    - valida firma/integridad,
    - ejecuta instalador de actualizacion.
5. El instalador ejecuta hooks de migracion de servicio:
    - stop servicio,
    - reemplazo binario,
    - reconfiguracion idempotente,
    - start servicio,
    - `health` check.
6. Electron relanza UI y reconecta con `kara-core`.

### Que validamos exactamente en GitHub Releases

- Publicacion de artefactos correctos por SO.
- Metadatos de update presentes (ejemplo `latest.yml` en ecosistema electron-builder).
- Firma de codigo y verificaciones de integridad.
- Capacidad de actualizar desde `vN` a `vN+1` sin romper el servicio.
- Notas de release con cambios de migracion cuando impacten al servicio.

### Politica de rollout recomendada

- Canal `beta` para validacion interna (early adopters).
- Canal `latest` para produccion estable.
- Rollout gradual (por porcentaje o por cohortes) antes de despliegue global.
- Mecanismo de rollback a release anterior si falla health check post-update.

### Hooks minimos que debe implementar el instalador

- Pre-update:
    - detectar servicio existente,
    - parar servicio con timeout,
    - backup opcional de config.
- Post-update:
    - registrar/actualizar servicio,
    - iniciar servicio,
    - validar `GET /health`,
    - marcar update como exitoso o revertir.

### Telemetria y auditoria de actualizaciones

Registrar eventos minimos:

- `update_available`
- `update_downloaded`
- `service_stop_started`
- `service_start_ok` o `service_start_failed`
- `health_check_ok` o `health_check_failed`
- `rollback_executed`

Con esto se puede auditar en produccion si las actualizaciones mantienen vivo el servicio.

### Validacion en CI/CD (obligatoria)

Antes de marcar un release como estable:

1. Instalar version anterior en runner/VM limpia.
2. Verificar servicio activo.
3. Ejecutar actualizacion desde GitHub Releases.
4. Verificar que el servicio sigue activo y responde `health`.
5. Verificar que la UI reconecta sin intervención manual.
6. Verificar desinstalacion limpia tras update.

# Criterios de aceptación

- Cerrar la ventana principal no finaliza la app.
- El icono aparece en el System Tray de Windows.
- Desde el tray se puede abrir, ocultar y cerrar realmente la app.
- La app puede arrancar con Windows.
- El atajo global muestra/oculta la mini ventana del asistente.
- La mini ventana funciona aunque la app principal esté oculta.
- Al bloquear Windows, las ventanas se ocultan.
- Al desbloquear Windows, los servicios internos siguen disponibles.
- La app no intenta mostrar UI sobre la pantalla de bloqueo.
- Todo IPC entre renderer y main está validado y limitado.
- No se expone Node.js directamente al renderer.
- El backend sigue activo como servicio aunque Electron no esté visible.
- Reinicios del backend son automáticos ante fallos.
- En Windows, la UI corre en sesión de usuario y el backend en servicio separado.
- La instalacion crea y arranca el servicio automaticamente sin pasos manuales.
- La actualizacion mantiene el servicio operativo y con reinicio automatico.
- La desinstalacion elimina el servicio y su registro del sistema.
- El flujo de actualizacion queda validado mediante GitHub Releases antes de promover a estable.
- Cada update ejecuta migracion de servicio + health check + rollback en caso de fallo.
