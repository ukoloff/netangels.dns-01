import { it, describe } from 'node:test'
import assert from 'assert'

it('Hello, world', x => {
  assert.equal(1, 1)
})

it.skip('Oops', async x => {
  assert.equal(1, 2)
})
