//
// NetAngels API
//
import { Resolver } from 'node:dns/promises'

const AUTH = "https://panel.netangels.ru/api/gateway/token/"
const API = 'https://api-ms.netangels.ru/api/v1/dns/'

const auth = doAuth()

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

export async function doAuth(key = process.env.NETANGELS_API_KEY) {
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
  return {
    headers: {
      authorization: `Bearer ${j.token}`
    }
  }
}

export async function zones() {
  let q = await fetch(`${API}zones`, await auth)
  return await q.json()
}

export async function RRs(zoneId) {
  let q = await fetch(`${API}zones/${zoneId}/records`, await auth)
  return await q.json()
}

export async function findRRs(name, type = null) {
  let result = []
  let zz = await zones()
  for (let z of zz.entities) {
    if ((0 == z.records_count) || (z.name != name && !name.endsWith('.' + z.name)))
      continue
    let Rs = await RRs(z.id)
    for (let r of Rs.entities) {
      if ((r.name != name) || (type && r.type != type))
        continue
      result.push(r)
    }
  }
  return result
}

export async function create(rec) {
  let params = await auth
  let q = await fetch(`${API}/records`, {
    ...params,
    method: 'POST',
    body: JSON.stringify(rec),
    headers: {
      'Content-Type': 'application/json',
      ...params.headers
    }
  })
  return await q.json()
}

export async function drop(rrId) {
  await fetch(`${API}records/${rrId}`, {
    ...await auth,
    method: 'DELETE',
  })
}
