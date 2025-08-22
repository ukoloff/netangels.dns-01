//
// Execute command from CLI
//
import { create, remove } from "./na.js"

export async function present(fqdn, text) {
  await create({
    name: fqdn,
    type: 'TXT',
    value: text,
  })
}

export async function cleanup(fqdn, text) {
  await remove(fqdn, {
    type: 'TXT',
    value: text,
  })
}
