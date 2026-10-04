---
category: cdk
type:
title: Overlay
subtitle: 浮层
cover:
---

通过函数创建浮层，并返回该浮层对应的实例，允许调用方在合适的范围内对浮层进行自定义。

## 何时使用

需要在某个触发元素附近展示浮层（例如 tooltip、popover 等）时使用。

- 通过函数可以创建浮层并返回该浮层对应的实例；
- 调用方可以通过函数返回的浮层实例，对视图进行绑定；
- 调用方可以通过返回的实例，控制浮层的初始化、显示、隐藏和销毁；
- 浮层的滚动策略：`close` 滚动时自动关闭浮层，`reposition` 滚动时重新定位（默认行为）；
- 浮层定位依赖于 `@popperjs/core`。

## API

### `useOverlay`

```ts
export const useOverlay: (options: OverlayOptions) => OverlayInstance
```

### `OverlayOptions`

| 属性 | 说明 | 类型 | 默认值 |
| --- | --- | --- | --- |
| `className` | 浮层容器的类名 | `string` | - |
| `visible` | 控制浮层显隐 | `boolean` | `false` |
| `scrollStrategy` | 滚动策略 | `'close' \| 'reposition'` | `'reposition'` |
| `disable` | 禁用浮层 | `boolean` | `false` |
| `arrowOffset` | 箭头与两端起始点的距离 | `number` | - |
| `showArrow` | 是否展示箭头 | `boolean` | `false` |
| `placement` | 浮层对齐方式 | `Placement` | `'bottom'` |
| `popperOptions` | popper 配置，优先级高于 `OverlayOptions` | `Partial<Options>` | - |
| `trigger` | 触发方式 | `'click' \| 'hover' \| 'focus'` | `'click'` |
| `allowEnter` | 是否允许鼠标进入浮层 | `boolean` | `true` |
| `offset` | 浮层偏移量 `[水平偏移, 垂直偏移]` | `[number, number]` | `[0, 0]` |
| `hideDelay` | 隐藏浮层的延迟，传 `false` 表示不需要延迟 | `number \| false` | `100` |
| `showDelay` | 显示浮层的延迟，传 `false` 表示不需要延迟 | `number \| false` | `0` |

### `OverlayInstance`

| 属性 | 说明 | 类型 |
| --- | --- | --- |
| `initialize` | 初始化浮层 | `() => void` |
| `show` | 显示浮层 | `() => void` |
| `hide` | 隐藏浮层 | `() => void` |
| `destroy` | 销毁浮层 | `() => void` |
| `overlayId` | 浮层的唯一 id | `string` |
| `visibility` | 当前浮层的显示状态，由 `visible` 和 `disable` 控制 | `ComputedRef<boolean>` |
| `overlayRef` | 浮层的真实 DOM 节点，需要绑定到视图上 | `Ref<RefElement>` |
| `update` | 更新浮层配置，未初始化时会先初始化 | `(options: Partial<OverlayOptions>) => void` |
| `arrowRef` | 箭头的真实 DOM 节点，`showArrow` 为 `false` 时不返回 | `Ref<RefElement>` |
| `triggerRef` | 触发元素的真实 DOM 节点，需要绑定到视图上 | `Ref<VueElement>` |
| `triggerEvents` | 需要手动绑定到触发元素上的事件 | `OverlayTriggerEvents` |
| `overlayEvents` | 需要手动绑定到浮层上的事件 | `OverlayPopperEvents` |
