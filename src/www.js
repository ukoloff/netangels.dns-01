//
// API server
//
import cluster from 'node:cluster'
import { createServer } from 'node:http'
import { json } from 'node:stream/consumers'
import watch from './watch.js'

export default function www(args) {
  if (cluster.isPrimary) {
    for (let i = 0; i < 3; i++)
      cluster.fork()
    cluster.on('exit', (worker, code, signal) => cluster.fork())
    watch()
    return
  }
  let srv = createServer(handler)
  srv.listen(80)
}

function handler(req, res) {
  let verb = req.url.replace(/^\/+|\/+$/g, '')
  let prsnt = 0
  switch (verb) {
    case 'alive':
      res.end('Ok')
      break
    case 'present': prsnt = 1
    case 'cleanup':
      res.end('+')
      break
    default:
      res.statusCode = 404
      res.end()
  }
}
