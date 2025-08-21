import { it, describe } from 'node:test'
import assert from 'node:assert'
import { Resolver } from 'node:dns/promises'

let _dns = resolver()

it('Test DNS', async _ => {
  let dns = await _dns
  let IPs = await dns.resolve('ekb.ru')
  assert.ok(IPs.length)
})

async function resolver(domain = 'ekb.ru') {
  const dns = new Resolver()
  dns.setServers(['8.8.8.8'])
  let ns = await dns.resolveNs(domain)
  let IPs = []
  for (let nserver of ns) {
    let ips = await dns.resolve4(nserver)
    IPs.push(...ips)
  }
  dns.setServers(IPs)
  return dns
}
