const { contextBridge } = require("electron");

contextBridge.exposeInMainWorld("kara", {
	version: "0.1.0",
	platform: process.platform,
	mode: process.env.NODE_ENV || "production",
});
