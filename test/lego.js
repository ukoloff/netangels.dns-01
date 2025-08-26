import fs from 'node:fs/promises'
import { X509Certificate } from 'node:crypto'
import { spawn } from 'node:child_process'
import { it, describe } from 'node:test'
import random from '../src/random.js'
import assert from 'node:assert'

describe('lego', $ => {
  it('uses exec interface', exec)
  it('uses http interface', http)
})

async function exec($) {
  await lego('exec', {
    EXEC_PATH: 'bin/na01',
    EXEC_POLLING_INTERVAL: 10,
    EXEC_PROPAGATION_TIMEOUT: 300,
  })
}

async function http($) {

}

function wait(child) {
  return new Promise(function (resolve, reject) {
    child
      .on('error', reject)
      .on('exit', resolve)
  })
}

async function lego(provider, env = {}) {
  process.env.LEGO_EMAIL = 'stas@ekb.ru'
  process.env.LEGO_SERVER = 'https://acme-staging-v02.api.letsencrypt.org/directory'
  for (const [k, v] of Object.entries(env)) {
    process.env[k] = v
  }

  let domain = `${await random()}.exec.uralhimmash.com`.toLowerCase()

  let child = spawn('lego', [
    '-a',
    '-dns', provider,
    '-d', domain,
    '--pfx',
    '--dns.resolvers', '8.8.8.8',
    'run'],
    { stdio: 'inherit' })
  let res = await wait(child)
  assert.equal(res, 0)
  let crt = new X509Certificate(await fs.readFile('./.lego/certificates/' + domain + '.crt'))
  assert.ok(crt.serialNumber)
  assert.equal(crt.subject, 'CN=' + domain)
  assert.equal(crt.subjectAltName, 'DNS:' + domain)

  for (let ext of 'key pfx crt json'.split(' ')) {
    assert.ok(await fs.stat('./.lego/certificates/' + domain + '.' + ext))
  }
}
