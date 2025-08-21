import { it, describe } from 'node:test'
import { doAuth, RRs, zones } from '../src/na.js'
import assert from 'node:assert'

describe('NetAngels', _ => {
  it('can authorize', async _ => {
    await doAuth()
  })

  it('enumerates zones', async _ => {
    let z = await zones()
    assert.ok(z.count)
  })

  it('enumerates records in zones', async _ => {
    let zz = await zones()
    for (let z of zz.entities) {
      if (z.is_technical_zone || 0 == z.records_count)
        continue
      let rs = await RRs(z.id)
      assert.ok(rs.count)
    }
  })

})
