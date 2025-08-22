import { setTimeout } from 'node:timers/promises'
import { it, describe } from 'node:test'
import { create, doAuth, drop, findRRs, remove, resolver, RRs, zones } from '../src/na.js'
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
    rs = await findRRs(`${await random()}.ekb.ru`)
    $.assert.strictEqual(rs.length, 0)
  })

  it('creates/deletes arbitrary records', async $ => {
    let r = await create({
      type: 'TXT',
      name: `${await random()}.test.uralhimmash.com`,
      value: `Hello, ${await random()}`,
      ttl: 300
    })
    let q = await drop(r.id)
    $.assert.strictEqual(q.type, 'TXT')
  })

  it('removes records by data', async $ => {
    let name = `${await random()}.test.uralhimmash.com`
    let rec = {
      type: 'TXT',
      value: `Bye, ${await random()}!`,
    }
    let r = await create({
      ...rec,
      name,
      ttl: 330,
    })
    let rs = await remove(name, rec)
    $.assert.deepEqual(1, rs.length)
  })

  it('measuring DNS delay', async $ => {
    for (let i = 0; i < 5; i++) {

      let name = `${await random()}.test.uralhimmash.com`
      let rec = {
        type: 'TXT',
        value: `Bye, ${await random()}!`,
      }
      let r = await create({
        ...rec,
        name,
        ttl: 330,
      })
      let rslvr = await resolver()
      let start = new Date()
      while (1) {
        try {
          let RRs = await rslvr.resolveTxt(name)
          $.assert.ok(RRs.length)
          break
        } catch (e) {
          process.stdout.write(`\r${Math.round((new Date() - start) / 1000)}...`)
          await setTimeout(3000)
        }
      }
      console.log('\rTime:', Math.round((new Date() - start) / 1000))
      let rs = await remove(name, rec)
      $.assert.deepEqual(1, rs.length)
    }

  })

})
