//
// API server
//
import cluster from 'node:cluster'
import { createServer } from 'node:http'
import { json } from 'node:stream/consumers'
import watch from './watch.js'
import { create, normalizeDomain, remove } from "./na.js"

export default function www(args) {
  if (cluster.isPrimary) {
    for (let i = 0; i < 3; i++)
      cluster.fork()
    cluster.on('exit', (worker, code, signal) => cluster.fork())
    if (process.env.DEV_WATCH)
      watch()
    return
  }
  let srv = createServer(handler)
  srv.listen(80)
}

async function handler(req, res) {
  try {
    let verb = req.url.replace(/^\/+|\/+$/g, '')
    let prsnt = 0
    switch (verb) {
      case 'alive':
        res.end('Ok')
        break
      case 'present': prsnt = 1
      case 'cleanup':
        let j = await json(req)
        let r
        if (prsnt) {
          r = await create({
            name: normalizeDomain(j.fqdn),
            type: 'TXT',
            value: j.value,
          })
        } else {
          r = await remove(normalizeDomain(j.fqdn), {
            type: 'TXT',
            value: j.value,
          })
        }
        res.setHeader('Content-Type', 'application/json')
        res.end(JSON.stringify(r))
        break
      default:
        res.statusCode = 404
        res.end()
    }
  } catch (e) {
    res.statusCode = 500
    res.setHeader('Content-Type', 'application/json')
    res.end(JSON.stringify({ error: e.message }))
  }
}
