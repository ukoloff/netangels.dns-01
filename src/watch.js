//
// Restart on file change
//
import { fstat, watch } from 'node:fs'
import cluster from 'node:cluster'

setImmediate(sleep)

export default function sentinel() {
  watch(import.meta.dirname, tick)
}

let state = 0

function sleep(ms = 100) {
  state = 1
  setTimeout($ => state = 0, ms)
}


function tick(event, name) {
  if (state) return
  sleep()
  console.log('Respawning...')
  for (var w of Object.values(cluster.workers))
    w.kill()
}
