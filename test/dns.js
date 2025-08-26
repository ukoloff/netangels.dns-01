import { it, describe } from 'node:test'
import { resolver } from '../src/na.js'

let _dns = resolver()

it.skip('Test DNS', async $ => {
  let dns = await _dns
  let IPs = await dns.resolve('ekb.ru')
  $.assert.ok(IPs.length)
})
