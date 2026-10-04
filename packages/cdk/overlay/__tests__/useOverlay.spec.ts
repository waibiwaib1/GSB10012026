import { defineComponent, h, onMounted, Teleport } from 'vue'
import { mount } from '@vue/test-utils'
import { useOverlay } from '../src/useOverlay'
import type { OverlayInstance, OverlayOptions } from '../src/types'

const defaultOptions: OverlayOptions = {
  visible: false,
  trigger: 'click',
  placement: 'bottom',
  scrollStrategy: 'reposition',
  offset: [0, 0],
  showDelay: false,
  hideDelay: false,
  allowEnter: true,
}

let overlay: OverlayInstance

const getComponent = (options: Partial<OverlayOptions> = {}) => {
  return defineComponent({
    setup() {
      overlay = useOverlay({ ...defaultOptions, ...options })
      onMounted(overlay.initialize)
      const { triggerRef, overlayRef, triggerEvents, overlayEvents } = overlay
      return () =>
        h('div', [
          h(
            'button',
            {
              ref: triggerRef,
              id: 'trigger',
              onClick: triggerEvents.onClick,
              onMouseenter: triggerEvents.onMouseEnter,
              onMouseleave: triggerEvents.onMouseLeave,
              onFocus: triggerEvents.onFocus,
              onBlur: triggerEvents.onBlur,
            },
            'trigger',
          ),
          h(
            Teleport as never,
            { to: 'body' },
            h(
              'div',
              {
                ref: overlayRef,
                id: 'overlay',
                onMouseenter: overlayEvents.onMouseEnter,
                onMouseleave: overlayEvents.onMouseLeave,
              },
              'overlay',
            ),
          ),
        ])
    },
  })
}

const getOverlay = () => document.body.querySelector('#overlay') as HTMLElement

describe('useOverlay.ts', () => {
  const setup = (options: Partial<OverlayOptions> = {}) => {
    const wrapper = mount(getComponent(options))
    return wrapper
  }

  beforeEach(() => {
    document.body.innerHTML = ''
  })

  test('initialize work', () => {
    setup()
    expect(overlay.overlayRef.value).not.toBeNull()
    expect(getOverlay().style.display).toBe('none')
  })

  test('overlayId work', () => {
    setup()
    const first = overlay.overlayId
    setup()
    expect(overlay.overlayId).not.toBe(first)
  })

  test('show and hide work', async () => {
    setup()
    overlay.show()
    expect(overlay.visibility.value).toBe(true)
    expect(getOverlay().style.display).toBe('block')

    overlay.hide()
    expect(overlay.visibility.value).toBe(false)
    expect(getOverlay().style.display).toBe('none')
  })

  test('visible option work', () => {
    setup({ visible: true })
    expect(overlay.visibility.value).toBe(true)
    expect(getOverlay().style.display).toBe('block')
  })

  test('disable work', async () => {
    setup({ visible: true, disable: true })
    expect(overlay.visibility.value).toBe(false)
    expect(getOverlay().style.display).toBe('none')
  })

  test('click trigger work', async () => {
    const wrapper = setup()
    const trigger = wrapper.find('#trigger')
    await trigger.trigger('click')
    expect(getOverlay().style.display).toBe('block')
    await trigger.trigger('click')
    expect(getOverlay().style.display).toBe('none')
  })

  test('hover trigger work', async () => {
    const wrapper = setup({ trigger: 'hover', showDelay: 10, hideDelay: 10 })
    const trigger = wrapper.find('#trigger')
    await trigger.trigger('mouseenter')
    expect(getOverlay().style.display).toBe('none')
    await new Promise(resolve => setTimeout(resolve, 20))
    expect(getOverlay().style.display).toBe('block')

    await trigger.trigger('mouseleave')
    await new Promise(resolve => setTimeout(resolve, 20))
    expect(getOverlay().style.display).toBe('none')
  })

  test('hover trigger without delay work', async () => {
    const wrapper = setup({ trigger: 'hover', showDelay: false, hideDelay: false })
    const trigger = wrapper.find('#trigger')
    await trigger.trigger('mouseenter')
    expect(getOverlay().style.display).toBe('block')
    await trigger.trigger('mouseleave')
    expect(getOverlay().style.display).toBe('none')
  })

  test('focus trigger work', async () => {
    const wrapper = setup({ trigger: 'focus' })
    const trigger = wrapper.find('#trigger')
    await trigger.trigger('focus')
    expect(getOverlay().style.display).toBe('block')
    await trigger.trigger('blur')
    expect(getOverlay().style.display).toBe('none')
  })

  test('allowEnter work', async () => {
    const wrapper = setup({ trigger: 'hover', hideDelay: 20 })
    const trigger = wrapper.find('#trigger')
    await trigger.trigger('mouseenter')
    await new Promise(resolve => setTimeout(resolve, 10))
    expect(getOverlay().style.display).toBe('block')

    await trigger.trigger('mouseleave')
    overlay.overlayEvents.onMouseEnter()
    await new Promise(resolve => setTimeout(resolve, 40))
    expect(getOverlay().style.display).toBe('block')

    overlay.overlayEvents.onMouseLeave()
    await new Promise(resolve => setTimeout(resolve, 40))
    expect(getOverlay().style.display).toBe('none')
  })

  test('scroll strategy close work', async () => {
    setup({ visible: true, scrollStrategy: 'close' })
    expect(getOverlay().style.display).toBe('block')
    window.dispatchEvent(new Event('scroll'))
    expect(getOverlay().style.display).toBe('none')
  })

  test('scroll strategy reposition work', async () => {
    setup({ visible: true, scrollStrategy: 'reposition' })
    expect(getOverlay().style.display).toBe('block')
    window.dispatchEvent(new Event('scroll'))
    expect(getOverlay().style.display).toBe('block')
  })

  test('update work', async () => {
    setup()
    overlay.update({ placement: 'top' })
    expect(getOverlay().style.display).toBe('none')
    overlay.update({ visible: true })
    expect(getOverlay().style.display).toBe('block')
  })

  test('destroy work', async () => {
    setup({ visible: true })
    expect(getOverlay().style.display).toBe('block')
    overlay.destroy()
    expect(getOverlay().style.display).toBe('none')
    overlay.show()
    expect(getOverlay().style.display).toBe('none')
  })

  test('arrow work', () => {
    setup({ showArrow: true })
    expect(overlay.arrowRef).toBeDefined()
  })
})
