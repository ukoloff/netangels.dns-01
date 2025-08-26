//
// NetAngels API
//
import { Resolver } from 'node:dns/promises'

const AUTH = "https://panel.netangels.ru/api/gateway/token/"
const API = 'https://api-ms.netangels.ru/api/v1/dns/'

export function normalizeDomain(domain) {
  return domain.replace(/[.]+$/, '').toLowerCase()
}

export async function resolver(domain = 'netangels.ru') {
  const dns = new Resolver({ timeout: 300, tries: 3 })
  // dns.setServers(['8.8.8.8'])
  let ns = await dns.resolveNs(domain)
  let IPs = []
  for (let nserver of ns) {
    let ips = await dns.resolve4(nserver)
    IPs.push(...ips)
  }
  dns.setServers(IPs)
  return dns
}

export async function auth(key = process.env.NETANGELS_API_KEY) {
  let params = new FormData()
  params.append('api_key', key)
  let auth = await fetch(AUTH, {
    method: 'POST',
    body: params,
  })
  let j = await auth.json()
  if (!j.token) {
    throw new Error('Failed to authorize to netangels.ru')
  }
  return j.token
}

let token

async function req(verb, options = {}) {
  token ||= auth()
  let q = await fetch(`${API}${verb}`, {
    ...options,
    headers: {
      ...options.headers,
      authorization: `Bearer ${await token}`
    }
  })
  if (q.ok)
    return await q.json()
}

export async function zones() {
  return await req('zones')
}

export async function RRs(zoneId) {
  return await req(`zones/${zoneId}/records`)
}

export async function findRRs(name, where = {}) {
  name = name.toLowerCase()
  let result = []
  let zz = await zones()
  for (let z of zz.entities) {
    if ((0 == z.records_count) || (z.name != name && !name.endsWith('.' + z.name)))
      continue
    let Rs = await RRs(z.id)
    RR: for (let r of Rs.entities) {
      if (r.name != name)
        continue
      for (const [key, value] of Object.entries(where)) {
        if (value === r[key])
          continue
        if (r.details && value === r.details[key])
          continue
        continue RR
      }
      result.push(r)
    }
  }
  return result
}

export async function create(rec) {
  return await req('records', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(rec),
  })
}

export async function drop(rrId) {
  return await req(`records/${rrId}`, { method: 'DELETE' })
}

export async function remove(name, where) {
  let rs = await findRRs(name, where)
  for (let r of rs) {
    await drop(r.id)
  }
  return rs
}
