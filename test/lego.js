import fs from 'node:fs/promises'
import { X509Certificate } from 'node:crypto'
import { spawn } from 'node:child_process'
import { it, describe } from 'node:test'
import random from '../src/random.js'

const Staging = 'https://acme-staging-v02.api.letsencrypt.org/directory'

describe('lego', lego)

function lego() {
  it('uses exec interface', exec)
  it('uses http interface', http)
}

async function exec($) {
  process.env.LEGO_EMAIL = 'stas@ekb.ru'
  process.env.EXEC_PATH = 'bin/na01'
  process.env.EXEC_PROPAGATION_TIMEOUT = 300
  process.env.LEGO_SERVER = Staging

  let domain = `${await random()}.exec.uralhimmash.com`.toLowerCase()

  let child = spawn('lego', [
    '-a',
    '-dns', 'exec',
    '-d', domain,
    '--pfx',
    '--dns.resolvers', '8.8.8.8',
    'run'],
    { stdio: 'inherit' })
  let res = await wait(child)
  $.assert.equal(res, 0)
  let crt = new X509Certificate(await fs.readFile('./.lego/certificates/' + domain + '.crt'))
  $.assert.ok(crt.serialNumber)
  $.assert.equal(crt.subject, 'CN=' + domain)
  $.assert.equal(crt.subjectAltName, 'DNS:' + domain)
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
