import help from './help.js'
import { present, cleanup } from './exec.js'

if (process.argv.length < 3)
  help()

let prsnt = false
switch (process.argv[2]) {
  case 'present': prsnt = true
  case 'cleanup':
    if (process.argv.length != 5)
      help()
    let args = process.argv.slice(3, 5)
    if (prsnt)
      await present(...args)
    else
      await cleanup(...args)
    break
  case 'www':
    if (process.argv.length != 3)
      help()
    break
  default:
    help()
}
