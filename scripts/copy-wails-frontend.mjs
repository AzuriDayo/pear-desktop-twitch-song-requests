import { cpSync, rmSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const dest = join(root, "cmd/main/frontend");
const src = join(root, "control-panel/build");

rmSync(dest, { recursive: true, force: true, maxRetries: process.platform === "win32" ? 10 : 0 });
cpSync(src, dest, { recursive: true });
