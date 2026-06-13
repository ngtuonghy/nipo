/**
 * Nipo Cloudflare Worker Backend
 * Manages CLI tunnel subdomains in KV storage and acts as a web gateway (via iframe)
 * to route incoming web traffic to active Cloudflare tunnels.
 */

import { Hono } from 'hono'

type Bindings = {
  NIPO_KV: KVNamespace
  API_SECRET: string
}

const app = new Hono<{ Bindings: Bindings }>()

/**
 * Register a new subdomain mapping in KV.
 * Used by the CLI client on startup.
 */
app.post('/register', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (authHeader !== `Bearer ${c.env.API_SECRET}`) {
    return c.json({ error: 'Unauthorized' }, 401)
  }
  try {
    const { subdomain, target } = await c.req.json()
    if (!subdomain || !target) {
      return c.json({ error: 'Missing subdomain or target' }, 400)
    }

    // Prevent overwriting an active subdomain
    const existing = await c.env.NIPO_KV.get(subdomain)
    if (existing) {
      return c.json({ error: 'Subdomain already in use. Please choose another one.' }, 403)
    }

    // Set 15-minute expiration (900 seconds)
    await c.env.NIPO_KV.put(subdomain, target, { expirationTtl: 900 })
    return c.json({ success: true, subdomain, target })
  } catch (err) {
    return c.json({ error: 'Bad Request' }, 400)
  }
})

/**
 * Renew the subdomain's TTL to prevent expiration.
 * CLI client sends this heartbeat periodically.
 */
app.post('/heartbeat', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (authHeader !== `Bearer ${c.env.API_SECRET}`) {
    return c.json({ error: 'Unauthorized' }, 401)
  }
  try {
    const { subdomain } = await c.req.json()
    if (!subdomain) return c.json({ error: 'Missing subdomain' }, 400)
    const target = await c.env.NIPO_KV.get(subdomain)
    if (!target) return c.json({ error: 'Subdomain not found' }, 404)

    // Extend expiration by 15 minutes
    await c.env.NIPO_KV.put(subdomain, target, { expirationTtl: 900 })
    return c.json({ success: true })
  } catch (err) {
    return c.json({ error: 'Bad Request' }, 400)
  }
})

/**
 * Unregister the subdomain immediately when the CLI client exits.
 */
app.delete('/unregister', async (c) => {
  const authHeader = c.req.header('Authorization')
  if (authHeader !== `Bearer ${c.env.API_SECRET}`) {
    return c.json({ error: 'Unauthorized' }, 401)
  }
  try {
    const { subdomain } = await c.req.json()
    if (!subdomain) return c.json({ error: 'Missing subdomain' }, 400)
    await c.env.NIPO_KV.delete(subdomain)
    return c.json({ success: true })
  } catch (err) {
    return c.json({ error: 'Bad Request' }, 400)
  }
})

/**
 * Wildcard handler that routes incoming HTTP requests to the target tunnel.
 */
app.all('*', async (c) => {
  const url = new URL(c.req.url)
  const hostname = url.hostname
  let subdomain = hostname.split('.')[0]

  // Support override for local testing
  if (url.searchParams.has('test_subdomain')) {
    subdomain = url.searchParams.get('test_subdomain')!
  }

  // Look up the tunnel target URL
  const targetUrlString = await c.env.NIPO_KV.get(subdomain)

  if (!targetUrlString) {
    return c.text('Nipo Tunnel not found or has expired.', 404)
  }

  // Map original path and query to the target URL
  const targetUrl = new URL(targetUrlString)
  targetUrl.pathname = url.pathname
  targetUrl.search = url.search

  // Proxy the request directly to the target tunnel
  return fetch(new Request(targetUrl.toString(), c.req.raw))
})

export default app
