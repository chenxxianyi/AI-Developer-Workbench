import { spawn, spawnSync } from 'node:child_process'
import net from 'node:net'

const port = await findAvailablePort(5173)
const serverUrl = `http://127.0.0.1:${port}`
const server = spawn(process.execPath, [
  './node_modules/vite/bin/vite.js',
  '--host', '127.0.0.1',
  '--port', String(port),
  '--strictPort',
], {
  cwd: process.cwd(),
  stdio: ['ignore', 'pipe', 'pipe'],
})

let stopped = false

server.stdout.on('data', (chunk) => {
  process.stdout.write(`[vite] ${chunk}`)
})

server.stderr.on('data', (chunk) => {
  process.stderr.write(`[vite] ${chunk}`)
})

function stopServer() {
  if (stopped || server.exitCode !== null || !server.pid) return
  stopped = true

  if (process.platform === 'win32') {
    spawnSync('taskkill', ['/pid', String(server.pid), '/T', '/F'], { stdio: 'ignore' })
    return
  }

  server.kill('SIGTERM')
}

process.on('exit', stopServer)
process.on('SIGINT', () => {
  stopServer()
  process.exit(130)
})
process.on('SIGTERM', () => {
  stopServer()
  process.exit(143)
})

try {
  await waitForServer(serverUrl, 30_000)

  const runner = spawn(process.execPath, ['./node_modules/playwright/cli.js', 'test', '--no-deps'], {
    cwd: process.cwd(),
    stdio: 'inherit',
    env: {
      ...process.env,
      PLAYWRIGHT_BASE_URL: serverUrl,
    },
  })

  const code = await new Promise((resolve) => {
    runner.on('exit', (exitCode) => resolve(exitCode ?? 1))
    runner.on('error', () => resolve(1))
  })

  stopServer()
  process.exit(Number(code))
} catch (error) {
  stopServer()
  console.error(error instanceof Error ? error.message : error)
  process.exit(1)
}

async function waitForServer(url, timeoutMs) {
  const startedAt = Date.now()
  let lastError

  while (Date.now() - startedAt < timeoutMs) {
    if (server.exitCode !== null) {
      throw new Error(`Vite dev server exited early with code ${server.exitCode}`)
    }

    try {
      const response = await fetch(url)
      if (response.ok) return
    } catch (error) {
      lastError = error
    }

    await new Promise((resolve) => setTimeout(resolve, 250))
  }

  throw new Error(`Timed out waiting for ${url}${lastError ? `: ${lastError}` : ''}`)
}

async function findAvailablePort(startPort) {
  for (let port = startPort; port < startPort + 100; port += 1) {
    if (await isPortAvailable(port)) return port
  }
  throw new Error(`No available port found between ${startPort} and ${startPort + 99}`)
}

function isPortAvailable(port) {
  return new Promise((resolve) => {
    const probe = net.createServer()
    probe.once('error', () => resolve(false))
    probe.once('listening', () => {
      probe.close(() => resolve(true))
    })
    probe.listen(port, '127.0.0.1')
  })
}
