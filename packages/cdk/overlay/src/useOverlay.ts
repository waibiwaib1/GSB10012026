import type { Instance as PopperInstance, Options as PopperOptions } from '@popperjs/core'
import type { Ref } from 'vue'

import { computed, reactive, ref, watch } from 'vue'
import { createPopper } from '@popperjs/core'
import { off, on } from '@idux/cdk/utils'

import type {
  OverlayInstance,
  OverlayOptions,
  OverlayPopperEvents,
  OverlayTriggerEvents,
  RefElement,
  VueElement,
} from './types'

const defaultOptions: OverlayOptions = {
  visible: false,
  scrollStrategy: 'reposition',
  disable: false,
  showArrow: false,
  placement: 'bottom',
  trigger: 'click',
  allowEnter: true,
  offset: [0, 0],
  hideDelay: 100,
  showDelay: 0,
}

let overlayId = 0

const resolveElement = (value: VueElement): RefElement => {
  if (!value) {
    return null
  }
  return (value as HTMLElement).nodeType === 1
    ? (value as HTMLElement)
    : ((value as { $el: RefElement }).$el as RefElement)
}

export function useOverlay(options: OverlayOptions): OverlayInstance {
  const state = reactive({ ...defaultOptions, ...options })

  const triggerRef = ref<VueElement>(null)
  const overlayRef = ref<RefElement>(null)
  const arrowRef = ref<RefElement>(null)

  const initialized = ref(false)
  const destroyed = ref(false)

  let popperInstance: PopperInstance | null = null
  let showTimer: number | null = null
  let hideTimer: number | null = null

  const visibility = computed(() => !!state.visible && !state.disable)

  const clearTimer = (timer: number | null) => {
    if (timer !== null) {
      clearTimeout(timer)
    }
  }

  const getPopperOptions = (): Partial<PopperOptions> => {
    const { placement, offset, showArrow, arrowOffset, popperOptions } = state
    const { modifiers: extraModifiers = [], ...restOptions } = popperOptions ?? {}
    const modifiers: NonNullable<PopperOptions['modifiers']> = [
      { name: 'offset', options: { offset } },
      ...extraModifiers,
    ]
    if (showArrow) {
      modifiers.push({
        name: 'arrow',
        options: { element: arrowRef.value, padding: arrowOffset },
      } as NonNullable<PopperOptions['modifiers']>[number])
    }
    return { placement, ...restOptions, modifiers }
  }

  const initPopper = () => {
    const triggerElement = resolveElement(triggerRef.value)
    const overlayElement = overlayRef.value
    if (!triggerElement || !overlayElement) {
      return
    }
    popperInstance = createPopper(triggerElement, overlayElement, getPopperOptions())
  }

  const scrollHandler = () => {
    if (state.scrollStrategy === 'close') {
      hide()
    } else {
      popperInstance?.update()
    }
  }

  const doShow = () => {
    if (!initialized.value || destroyed.value) {
      return
    }
    const overlayElement = overlayRef.value
    if (!overlayElement) {
      return
    }
    if (!popperInstance) {
      initPopper()
    }
    overlayElement.style.display = 'block'
    void popperInstance?.update()
    on(window, 'scroll', scrollHandler)
  }

  const doHide = () => {
    const overlayElement = overlayRef.value
    if (overlayElement) {
      overlayElement.style.display = 'none'
    }
    off(window, 'scroll', scrollHandler)
  }

  watch(
    visibility,
    value => {
      if (value) {
        doShow()
      } else {
        doHide()
      }
    },
    { flush: 'sync' },
  )

  const initialize = () => {
    if (initialized.value || destroyed.value) {
      return
    }
    initialized.value = true
    if (visibility.value) {
      doShow()
    } else {
      doHide()
    }
  }

  const show = () => {
    if (destroyed.value) {
      return
    }
    state.visible = true
  }

  const hide = () => {
    clearTimer(showTimer)
    state.visible = false
  }

  const destroy = () => {
    clearTimer(showTimer)
    clearTimer(hideTimer)
    state.visible = false
    off(window, 'scroll', scrollHandler)
    popperInstance?.destroy()
    popperInstance = null
    destroyed.value = true
  }

  const update = (options: Partial<OverlayOptions>) => {
    Object.assign(state, options)
    if (!initialized.value) {
      initialize()
    } else {
      popperInstance?.setOptions(getPopperOptions())
    }
  }

  const showWithDelay = () => {
    clearTimer(hideTimer)
    if (state.showDelay === false) {
      show()
      return
    }
    clearTimer(showTimer)
    showTimer = setTimeout(show, state.showDelay) as unknown as number
  }

  const hideWithDelay = () => {
    clearTimer(showTimer)
    if (state.hideDelay === false) {
      hide()
      return
    }
    clearTimer(hideTimer)
    hideTimer = setTimeout(hide, state.hideDelay) as unknown as number
  }

  const triggerEvents: OverlayTriggerEvents = {
    onClick: () => {
      if (state.trigger === 'click') {
        visibility.value ? hide() : show()
      }
    },
    onMouseEnter: () => {
      if (state.trigger === 'hover') {
        showWithDelay()
      }
    },
    onMouseLeave: () => {
      if (state.trigger === 'hover') {
        hideWithDelay()
      }
    },
    onFocus: () => {
      if (state.trigger === 'focus') {
        show()
      }
    },
    onBlur: () => {
      if (state.trigger === 'focus') {
        hide()
      }
    },
  }

  const overlayEvents: OverlayPopperEvents = {
    onMouseEnter: () => {
      if (state.allowEnter) {
        clearTimer(hideTimer)
      }
    },
    onMouseLeave: () => {
      if (state.trigger === 'hover') {
        hideWithDelay()
      }
    },
  }

  const instance: OverlayInstance = {
    initialize,
    show,
    hide,
    destroy,
    overlayId: `ix-overlay-${overlayId++}`,
    visibility,
    overlayRef: overlayRef as Ref<RefElement>,
    update,
    triggerRef: triggerRef as Ref<VueElement>,
    triggerEvents,
    overlayEvents,
  }

  if (state.showArrow) {
    instance.arrowRef = arrowRef as Ref<RefElement>
  }

  return instance
}
