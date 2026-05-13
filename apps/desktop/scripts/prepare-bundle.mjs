import {
	copyFileSync,
	cpSync,
	existsSync,
	mkdirSync,
	readdirSync,
	rmSync,
} from "node:fs";
import path from "node:path";
import process from "node:process";

const desktopRoot = process.cwd();
const webOutDir = path.resolve(desktopRoot, "../web/out");
const coreDistDir = path.resolve(desktopRoot, "../core/dist");
const desktopWebDistDir = path.resolve(desktopRoot, "web-dist");
const desktopCoreBinDir = path.resolve(desktopRoot, "core-bin");

function resetDir(dirPath) {
	rmSync(dirPath, { recursive: true, force: true });
	mkdirSync(dirPath, { recursive: true });
}

function requirePath(pathToCheck, hint) {
	if (!existsSync(pathToCheck)) {
		throw new Error(`${pathToCheck} does not exist. ${hint}`);
	}
}

function copyWebBundle() {
	requirePath(
		webOutDir,
		"Run npm --workspace web run build before packaging desktop.",
	);
	resetDir(desktopWebDistDir);
	cpSync(webOutDir, desktopWebDistDir, { recursive: true });
}

function resolveCoreBinaryName() {
	const names = ["kara-core", "kara-core.exe"];
	for (const name of names) {
		if (existsSync(path.join(coreDistDir, name))) {
			return name;
		}
	}

	const candidates = existsSync(coreDistDir) ? readdirSync(coreDistDir) : [];
	throw new Error(
		`No core binary found in ${coreDistDir}. Found: ${candidates.join(", ") || "none"}`,
	);
}

function copyCoreBinary() {
	requirePath(
		coreDistDir,
		"Run npm --workspace core run build before packaging desktop.",
	);
	resetDir(desktopCoreBinDir);

	const binaryName = resolveCoreBinaryName();
	const fromPath = path.join(coreDistDir, binaryName);
	const targetName =
		process.platform === "win32" ? "kara-core.exe" : "kara-core";
	copyFileSync(fromPath, path.join(desktopCoreBinDir, targetName));
}

copyWebBundle();
copyCoreBinary();
console.log("Desktop bundle prepared: web-dist + core-bin");
