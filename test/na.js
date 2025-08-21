import { it, describe } from 'node:test'
import { doAuth, zones } from '../src/na.js'
import assert from 'node:assert'

describe('NetAngels', _ => {
  it('can authorize', async _ => {
    await doAuth()
  })

  it('enumerates zones', async _ => {
    let z = await zones()
    assert.ok(z.count)
  })
})
