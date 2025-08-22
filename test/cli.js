import { spawn } from 'node:child_process'
import { it, describe } from 'node:test'
import { create, doAuth, drop, findRRs, remove, resolver, RRs, zones } from '../src/na.js'
import random from '../src/random.js'

describe('CLI interface', $ => {

  it('creates TXT RRs', async $ => {
    let name = `${await random()}.cli.uralhimmash.com`
    let value = `Hi, ${await random()}!`
    let child = spawn('node', ['.', 'present', name, value], {stdio: 'inherit'})

  })

  it('removes TXT RRs', async $ => {
    let name = `${random()}.cli.uralhimmash.com`
    let value = `Oops, ${random()}!`

  })

})
