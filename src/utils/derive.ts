import { proxy, subscribe, getVersion } from '../vanilla'

type DeriveGet = <T extends object>(proxyObject: T) => T

type Subscriptions<U extends object> = Map<
  object,
  [callbackMap: Map<keyof U, () => void>, unsubscribe: () => void]
>

const subscriptionsCache = new WeakMap<object, Subscriptions<object>>()

const getSubscriptions = (proxyObject: object) => {
  let subscriptions = subscriptionsCache.get(proxyObject)
  if (!subscriptions) {
    subscriptions = new Map()
    subscriptionsCache.set(proxyObject, subscriptions)
  }
  return subscriptions
}

export const unstable_getDeriveSubscriptions = getSubscriptions

/**
 * derive
 *
 * This creates derived properties and attaches them
 * to a new proxy object or an existing proxy object.
 *
 * @example
 * import { proxy } from 'valtio'
 * import { derive } from 'valtio/utils'
 *
 * const state = proxy({
 *   count: 1,
 * })
 *
 * const derivedState = derive({
 *   doubled: (get) => get(state).count * 2,
 * })
 *
 * derive({
 *   tripled: (get) => get(state).count * 3,
 * }, {
 *   proxy: state,
 * })
 */
export const derive = <T extends object, U extends object>(
  derivedFns: {
    [K in keyof U]: (get: DeriveGet) => U[K]
  },
  options?: {
    proxy?: T
    sync?: boolean
  }
) => {
  const proxyObject = (options?.proxy || proxy({})) as U
  const notifyInSync = options?.sync
  const subscriptions: Subscriptions<U> = getSubscriptions(proxyObject)

  const addSubscription = (
    p: object,
    key: keyof U,
    callback: () => void
  ) => {
    const subscription = subscriptions.get(p)
    if (subscription) {
      subscription[0].set(key, callback)
    } else {
      const unsubscribe = subscribe(
        p,
        (ops) => {
          if (
            p === proxyObject &&
            ops.every((op) => op[1].length === 1 && op[1][0] in derivedFns)
          ) {
            return
          }
          subscriptions.get(p)?.[0].forEach((cb) => {
            cb()
          })
        },
        notifyInSync
      )
      subscriptions.set(p, [new Map([[key, callback]]), unsubscribe])
    }
  }

  const removeSubscription = (p: object, key: keyof U) => {
    const subscription = subscriptions.get(p)
    if (subscription) {
      const [callbackMap, unsubscribe] = subscription
      callbackMap.delete(key)
      if (!callbackMap.size) {
        unsubscribe()
        subscriptions.delete(p)
      }
    }
  }

  ;(Object.keys(derivedFns) as (keyof U)[]).forEach((key) => {
    if (Object.getOwnPropertyDescriptor(proxyObject, key)) {
      throw new Error('object property already defined')
    }
    const fn = derivedFns[key]
    let lastDependencies: Map<object, number> | null = null

    const evaluate = () => {
      if (lastDependencies) {
        if (
          Array.from(lastDependencies).every(
            ([p, version]) => getVersion(p) === version
          )
        ) {
          return
        }
      }

      const dependencies = new Map<object, number>()
      const get = <P extends object>(p: P) => {
        dependencies.set(p, getVersion(p))
        return p
      }
      const value = fn(get)

      const updateSubscriptions = () => {
        dependencies.forEach((_, p) => {
          if (!lastDependencies?.has(p)) {
            addSubscription(p, key, evaluate)
          }
        })
        lastDependencies?.forEach((_, p) => {
          if (!dependencies.has(p)) {
            removeSubscription(p, key)
          }
        })
        lastDependencies = dependencies
      }

      if (value instanceof Promise) {
        value.then(() => {
          updateSubscriptions()
          evaluate()
        })
      } else {
        updateSubscriptions()
      }
      proxyObject[key] = value
    }

    evaluate()
  })

  return proxyObject as T & U
}

/**
 * underive
 *
 * This stops derived properties from evaluating. It stops all subscriptions,
 * or only subscriptions for the specified keys. It can also delete derived
 * properties so new ones can be attached.
 *
 * @example
 * import { proxy } from 'valtio'
 * import { derive, underive } from 'valtio/utils'
 *
 * const state = proxy({ count: 1 })
 * const derivedState = derive({
 *   doubled: (get) => get(state).count * 2,
 * })
 *
 * underive(derivedState)
 */
export const underive = <T extends object, U extends object>(
  proxyObject: T & U,
  options?: {
    delete?: boolean
    keys?: (keyof U)[]
  }
) => {
  const subscriptions: Subscriptions<U> = getSubscriptions(proxyObject)
  const keysToDelete = options?.delete ? new Set<keyof U>() : null

  subscriptions.forEach(([callbackMap, unsubscribe], p) => {
    if (options?.keys) {
      options.keys.forEach((key) => {
        if (callbackMap.has(key)) {
          callbackMap.delete(key)
          if (keysToDelete) {
            keysToDelete.add(key)
          }
        }
      })
    } else {
      if (keysToDelete) {
        Array.from(callbackMap.keys()).forEach((key) => {
          keysToDelete.add(key)
        })
      }
      callbackMap.clear()
    }

    if (!callbackMap.size) {
      unsubscribe()
      subscriptions.delete(p)
    }
  })

  if (keysToDelete) {
    keysToDelete.forEach((key) => {
      delete proxyObject[key]
    })
  }
}
