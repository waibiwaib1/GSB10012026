/**
 * @author Yosuke Ota
 */
'use strict'

const RuleTester = require('eslint').RuleTester
const rule = require('../../../lib/rules/no-restricted-component-names')

const tester = new RuleTester({
  parser: require.resolve('vue-eslint-parser'),
  parserOptions: { ecmaVersion: 2020, sourceType: 'module' }
})

tester.run('no-restricted-component-names', rule, {
  valid: [
    {
      filename: 'test.vue',
      code: `
      <script>
      export default {
        name: 'Allow'
      }
      </script>
      `,
      options: ['Foo', '/^Baz(One|Two)/']
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      export default {
        name: 'Allow',
        components: {
          Allow: {}
        }
      }
      </script>
      `,
      options: ['Foo', '/^Baz(One|Two)/']
    },
    {
      filename: 'test.vue',
      code: `
      <script setup>
      defineOptions({ name: 'Allow' })
      </script>
      `,
      options: ['Foo', '/^Baz(One|Two)/']
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      export default {
        name: 'Foo'
      }
      </script>
      `
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      import { defineComponent } from 'vue'
      export default defineComponent({
        name: 'Allow'
      })
      </script>
      `,
      options: ['Foo']
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      const app = Vue.component('Allow', {})
      </script>
      `,
      options: ['Foo']
    }
  ],
  invalid: [
    {
      filename: 'test.vue',
      code: `
      <script>
      export default {
        name: 'Foo'
      }
      </script>
      `,
      options: ['Foo', '/^Baz(One|Two)/'],
      errors: [
        {
          message: 'Using `Foo` name is not allowed.',
          line: 4
        }
      ]
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      export default {
        name: 'BazTwo'
      }
      </script>
      `,
      options: ['Foo', '/^Baz(One|Two)/'],
      errors: [
        {
          message: 'Using `BazTwo` name is not allowed.',
          line: 4
        }
      ]
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      export default {
        components: {
          Foo: {},
          BazOne: {}
        }
      }
      </script>
      `,
      options: ['Foo', '/^Baz(One|Two)/'],
      errors: [
        {
          message: 'Using `Foo` name is not allowed.',
          line: 5
        },
        {
          message: 'Using `BazOne` name is not allowed.',
          line: 6
        }
      ]
    },
    {
      filename: 'test.vue',
      code: `
      <script setup>
      defineOptions({ name: 'Foo' })
      </script>
      `,
      options: [{ name: 'Foo', message: 'Custom message.' }],
      errors: [
        {
          message: 'Custom message.',
          line: 3
        }
      ]
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      const app = Vue.component('Foo', {})
      </script>
      `,
      options: ['Foo'],
      errors: [
        {
          message: 'Using `Foo` name is not allowed.',
          line: 3
        }
      ]
    },
    {
      filename: 'test.vue',
      code: `
      <script>
      import { defineComponent } from 'vue'
      export default defineComponent({
        name: \`Foo\`
      })
      </script>
      `,
      options: ['Foo'],
      errors: [
        {
          message: 'Using `Foo` name is not allowed.',
          line: 5
        }
      ]
    }
  ]
})
