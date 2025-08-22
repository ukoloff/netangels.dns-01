import { spawn } from 'node:child_process'
import { it, describe, before, after } from 'node:test'
import { create, remove } from '../src/na.js'
import random from '../src/random.js'
import { setTimeout } from 'node:timers/promises'

describe('Web API', $ => {
  let server

  before(async $ => {
    server = spawn('node', ['.', 'www'], {stdio: 'inherit'})
    await new Promise((resolve, reject) => {
      server
        .on('spawn', resolve)
        .on('error', reject)
    })
    while (1) {
      try {
        await fetch('http://localhost/alive')
        break
      } catch(e) {
        await setTimeout(100)
      }
    }
  })

  after($ => {
    server.kill()
  })

  it('creates TXT RRs', async $ => {
    let q = await fetch('http://localhost/present')
    let t = await q.text()
    $.assert.ok(t)
    return
    let name = `${await random()}.cli.uralhimmash.com`
    let value = `Hi, ${await random()}!`
    let child = spawn('node', ['.', 'present', name, value], { stdio: 'inherit' })
    let res = await wait(child)
    $.assert.equal(res, 0)
    let RRs = await remove(name, {
      type: 'TXT',
      value,
    })
    $.assert.equal(RRs.length, 1)
  })

  it('removes TXT RRs', async $ => {
    let q = await fetch('http://localhost/cleanup')
    let t = await q.text()
    $.assert.ok(t)
    return
    let name = `${await random()}.cli.uralhimmash.com`
    let value = `Oops, ${random()}!`
    let r = await create({
      name,
      type: 'TXT',
      value,
      ttl: 301,
    })
    $.assert.equal(r.type, 'TXT')
    let child = spawn('node', ['.', 'cleanup', name, value], { stdio: 'inherit' })
    let res = await wait(child)
    $.assert.equal(res, 0)
    let RRs = await remove(name, {
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
