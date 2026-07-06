import { execSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
execSync("vp run build", { cwd: root, stdio: "inherit" });
execSync("node scripts/copy-wails-frontend.mjs", { cwd: root, stdio: "inherit" });
