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
