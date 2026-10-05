---
pageClass: rule-details
sidebarDepth: 0
title: vue/no-restricted-component-names
description: disallow specific component names
---
# vue/no-restricted-component-names

> disallow specific component names

- :exclamation: <badge text="This rule has not been released yet." vertical="middle" type="error"> ***This rule has not been released yet.*** </badge>

## :book: Rule Details

This rule allows you to specify component names that you don't want to use in your application.

## :wrench: Options

This rule takes a list of strings, where each string is a component name or pattern to be restricted:

```json
{
  "vue/no-restricted-component-names": ["error", "Foo", "/^Baz(One|Two)/"]
}
```

<eslint-code-block :rules="{'vue/no-restricted-component-names': ['error', 'Foo', '/^Baz(One|Two)/']}">

```vue
<script>
/* ✗ BAD */
export default {
  name: 'Foo'
}
</script>
```

</eslint-code-block>

<eslint-code-block :rules="{'vue/no-restricted-component-names': ['error', 'Foo', '/^Baz(One|Two)/']}">

```vue
<script>
/* ✗ BAD */
export default {
  name: 'BazTwo'
}
</script>
```

</eslint-code-block>

<eslint-code-block :rules="{'vue/no-restricted-component-names': ['error', 'Foo', '/^Baz(One|Two)/']}">

```vue
<script>
/* ✓ GOOD */
export default {
  name: 'Allow'
}
</script>
```

</eslint-code-block>

Alternatively, the rule also accepts objects, where the name and a custom message can be specified:

```json
{
  "vue/no-restricted-component-names": [
    "error",
    {
      "name": "Foo",
      "message": "Use \"Bar\" instead of \"Foo\"."
    }
  ]
}
```

## :couple: Related Rules

- [vue/no-reserved-component-names](./no-reserved-component-names.md)

## :mag: Implementation

- [Rule source](https://github.com/vuejs/eslint-plugin-vue/blob/master/lib/rules/no-restricted-component-names.js)
- [Test source](https://github.com/vuejs/eslint-plugin-vue/blob/master/tests/lib/rules/no-restricted-component-names.js)
