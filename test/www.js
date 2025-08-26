import { spawn } from 'node:child_process'
import { it, describe, before, after } from 'node:test'
import { create, remove } from '../src/na.js'
import random from '../src/random.js'
import { setTimeout } from 'node:timers/promises'

describe.skip('Web API', $ => {
  let server

  before(async $ => {
    server = spawn('node', ['.', 'www'], { stdio: 'inherit' })
    await new Promise((resolve, reject) => {
      server
        .on('spawn', resolve)
        .on('error', reject)
    })
    let t
    while (1) {
      try {
        let q = await fetch('http://localhost/alive')
        t = await q.text()
        break
      } catch (e) {
        await setTimeout(100)
      }
    }
  })

  after(async $ => {
    await setTimeout(300)
    server.kill()
  })

  it('creates TXT RRs', async $ => {
    let fqdn = `${await random()}.web.uralhimmash.com`
    let value = `Oh, ${await random()}!`

    let q = await fetch('http://localhost/present', {
      method: 'POST',
      body: JSON.stringify({ fqdn, value }),
      headers: { 'Content-Type': 'application/json' }
    })
    let t = await q.json()
    $.assert.equal(t.type, 'TXT')

    let RRs = await remove(fqdn, {
      type: 'TXT',
      value,
    })
    $.assert.equal(RRs.length, 1)
  })

  it('removes TXT RRs', async $ => {
    let fqdn = `${await random()}.web.uralhimmash.com`
    let value = `Wow, ${await random()}!`

    let r = await create({
      name: fqdn,
      type: 'TXT',
      value,
      ttl: 301,
    })
    $.assert.equal(r.type, 'TXT')

    let q = await fetch('http://localhost/cleanup', {
      method: 'POST',
      body: JSON.stringify({ fqdn, value }),
      headers: { 'Content-Type': 'application/json' }
    })
    let t = await q.json()
    $.assert.ok(t.length, 1)

    let RRs = await remove(fqdn, {
      type: 'TXT',
      value,
    })
    $.assert.equal(RRs.length, 0)
  })

})


function wait(child) {
  return new Promise(function (resolve, reject) {
    child
      .on('error', reject)
      .on('exit', resolve)
  })
}
