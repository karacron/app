const { app, BrowserWindow } = require("electron");
const path = require("node:path");

function resolveWebEntry() {
	if (!app.isPackaged) {
		return path.join(__dirname, "renderer", "index.html");
	}

	return path.join(__dirname, "web-dist", "index.html");
}

function resolveCoreBinary() {
	const coreName =
		process.platform === "win32" ? "kara-core.exe" : "kara-core";
	if (!app.isPackaged) {
		return path.join(__dirname, "core-bin", coreName);
	}

	return path.join(process.resourcesPath, "core-bin", coreName);
}

function createWindow() {
	const win = new BrowserWindow({
		width: 1200,
		height: 800,
		minWidth: 980,
		minHeight: 680,
		webPreferences: {
			preload: path.join(__dirname, "preload.js"),
			contextIsolation: true,
			nodeIntegration: false,
		},
	});

	win.loadFile(resolveWebEntry());
}

app.whenReady().then(() => {
	createWindow();

	app.on("activate", () => {
		if (BrowserWindow.getAllWindows().length === 0) {
			createWindow();
		}
	});
});

app.on("window-all-closed", () => {
	if (process.platform !== "darwin") {
		app.quit();
	}
});

module.exports = {
	resolveCoreBinary,
};
