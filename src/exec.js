//
// Execute command from CLI
//
import { create, normalizeDomain, remove } from "./na.js"

export async function present(fqdn, text) {
  await create({
    name:  normalizeDomain(fqdn),
    type: 'TXT',
    value: text,
  })
}

export async function cleanup(fqdn, text) {
  await remove(normalizeDomain(fqdn), {
    type: 'TXT',
    value: text,
  })
}
