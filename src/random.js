//
// Generate random string
//
import { randomBytes } from 'node:crypto'

export default async function random(n = 9) {
  return new Promise(run)

  function run(resolve, reject) {
    randomBytes(n, cb)

    function cb(err, buf) {
      if (err) {
        reject(err)
        return
      }
      resolve(buf.toString('base64').replace(/\W+/g, ''))
    }
  }
}
