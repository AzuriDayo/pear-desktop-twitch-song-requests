Set-Location "$(git rev-parse --show-toplevel)"
pnpm run -r build
node -e "require('fs').rmSync('./cmd/main/build', { recursive: true, force: true, maxRetries: process.platform === 'win32' ? 10 : 0 })"
Copy-Item -r control-panel/build cmd/main
