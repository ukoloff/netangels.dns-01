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
  console.log(req)
  res.end('Hoo!')
}
