//
// Print usage
//
import { basename } from 'node:path'

export default function help() {
  let bname = basename(process.argv[1])
  console.info(`Usage:
\tnode ${bname} present fqdn text\t- Add DNS TXT record
\tnode ${bname} cleanup fqdn text\t- Remove DNS TXT record
\tnode ${bname} www\t\t\t- Start WWW server`)
  process.exit(1)
}
