import { it, describe } from 'node:test'
import { auth } from '../src/na.js'

describe('NetAngels', _ => {
  it('can authorize', async _ => {
    let h = await auth()
    return h
  })
})
