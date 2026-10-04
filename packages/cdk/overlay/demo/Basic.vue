<template>
  <ix-button ref="triggerRef" @click="triggerEvents.onClick">Click</ix-button>
  <teleport to="body">
    <div ref="overlayRef" @mouseenter="overlayEvents.onMouseEnter()" @mouseleave="overlayEvents.onMouseLeave()">
      tooltip
    </div>
  </teleport>
</template>

<script lang="ts">
import { defineComponent, onMounted } from 'vue'
import { useOverlay } from '@idux/cdk/overlay'

export default defineComponent({
  name: 'Basic',
  setup() {
    const { initialize, overlayRef, triggerRef, triggerEvents, overlayEvents } = useOverlay({
      visible: false,
      trigger: 'click',
      placement: 'bottom',
      offset: [0, 0],
      showDelay: false,
      hideDelay: 1000,
      allowEnter: true,
    })
    onMounted(initialize)

    return { overlayRef, triggerRef, triggerEvents, overlayEvents }
  },
})
</script>
