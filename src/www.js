//
// API server
//
import cluster from 'node:cluster'
import { createServer } from 'node:http'
import { json } from 'node:stream/consumers'

export default function www(args) {
  if (cluster.isPrimary) {
    for (let i = 0; i < 3; i++)
      cluster.fork()
    cluster.on('exit', (worker, code, signal) => cluster.fork())
    return
  }
  let srv = createServer(handler)
  srv.listen(80)
}

function handler(req, res) {
  console.log(req)
  res.end('Hi!')
}
