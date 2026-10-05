import { proxy, snapshot, subscribe } from '../src/index'
import { derive, underive } from '../src/utils'

it('creates a derived proxy', async () => {
  const state = proxy({ count: 0, text: '' })
  const derived = derive({
    doubled: (get) => get(state).count * 2,
  })

  expect(snapshot(derived)).toMatchObject({ doubled: 0 })

  state.count += 1
  await Promise.resolve()
  expect(snapshot(derived)).toMatchObject({ doubled: 2 })

  state.text = 'hello'
  await Promise.resolve()
  expect(snapshot(derived)).toMatchObject({ doubled: 2 })
})

it('attaches derived properties to an existing proxy', async () => {
  const state = proxy({ count: 0 })
  derive(
    {
      doubled: (get) => get(state).count * 2,
    },
    { proxy: state }
  )

  state.count += 1
  await Promise.resolve()
  expect(snapshot(state)).toMatchObject({ count: 1, doubled: 2 })
})

it('derives synchronously when sync is enabled', () => {
  const state = proxy({ count: 0 })
  const callback = jest.fn()
  const derived = derive(
    {
      doubled: (get) => get(state).count * 2,
    },
    { sync: true }
  )
  subscribe(derived, callback, true)

  state.count += 1

  expect(snapshot(derived)).toMatchObject({ doubled: 2 })
  expect(callback).toBeCalledTimes(1)
})

it('stops derived updates with underive', async () => {
  const state = proxy({ count: 0 })
  const derived = derive({
    doubled: (get) => get(state).count * 2,
  })

  underive(derived)
  state.count += 1
  await Promise.resolve()

  expect(snapshot(derived)).toMatchObject({ doubled: 0 })
})
