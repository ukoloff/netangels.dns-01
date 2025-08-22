import { it, describe } from 'node:test'
import { create, doAuth, drop, findRRs, RRs, zones } from '../src/na.js'
import random from '../src/random.js'

describe('NetAngels', _ => {
  it('can authorize', async _ => {
    await doAuth()
  })

  it('enumerates zones', async $ => {
    let z = await zones()
    $.assert.ok(z.count)
  })

  it('enumerates records in zones', async $ => {
    let zz = await zones()
    for (let z of zz.entities) {
      if (z.is_technical_zone || 0 == z.records_count)
        continue
      let rs = await RRs(z.id)
      $.assert.ok(rs.count)
    }
  })

  it('finds records', async $ => {
    let rs = await findRRs('ekb.ru')
    $.assert.ok(rs.length)
    let xs = await findRRs(`${await random()}.ekb.ru`)
    $.assert.equal(xs.length, 0)
  })

  it('creates/deletes arbitrary records', async $ => {
    let r = await create({
      type: 'TXT',
      name: `${await random()}.test.uralhimmash.com`,
      value: `Hello, ${await random()}`,
      ttl: 300
    })
    await drop(r.id)
  })

})
