import { it, describe } from 'node:test'
import { doAuth } from '../src/na.js'

describe('NetAngels', _ => {
  it('can authorize', async _ => {
    await doAuth()
  })
})
