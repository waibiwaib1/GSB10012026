/**
 * @author Yosuke Ota
 * See LICENSE file in root directory for full license.
 */
'use strict'

const utils = require('../utils')
const regexp = require('../utils/regexp')

/**
 * @typedef {object} ParsedOption
 * @property { (name: string) => boolean } test
 * @property {string|undefined} [message]
 */

/**
 * @param {string} str
 * @returns {(str: string) => boolean}
 */
function buildMatcher(str) {
  if (regexp.isRegExp(str)) {
    const re = regexp.toRegExp(str)
    return (s) => {
      re.lastIndex = 0
      return re.test(s)
    }
  }
  return (s) => s === str
}

/**
 * @param {string|{name:string, message?: string}} option
 * @returns {ParsedOption}
 */
function parseOption(option) {
  if (typeof option === 'string') {
    const matcher = buildMatcher(option)
    return {
      test(name) {
        return matcher(name)
      }
    }
  }
  const parsed = parseOption(option.name)
  parsed.message = option.message
  return parsed
}

/**
 * @param {Expression | SpreadElement} node
 * @returns {node is (Literal | TemplateLiteral)}
 */
function canVerify(node) {
  return (
    node.type === 'Literal' ||
    (node.type === 'TemplateLiteral' &&
      node.expressions.length === 0 &&
      node.quasis.length === 1)
  )
}

module.exports = {
  meta: {
    type: 'problem',
    docs: {
      description: 'disallow specific component names',
      categories: undefined,
      url: 'https://eslint.vuejs.org/rules/no-restricted-component-names.html'
    },
    fixable: null,
    schema: {
      type: 'array',
      items: {
        oneOf: [
          { type: ['string'] },
          {
            type: 'object',
            properties: {
              name: { type: 'string' },
              message: { type: 'string', minLength: 1 }
            },
            required: ['name'],
            additionalProperties: false
          }
        ]
      },
      uniqueItems: true,
      minItems: 0
    },
    messages: {
      // eslint-disable-next-line eslint-plugin/report-message-format
      restrictedName: '{{message}}'
    }
  },
  /** @param {RuleContext} context */
  create(context) {
    /** @type {ParsedOption[]} */
    const options = context.options.map(parseOption)

    /**
     * @param {ESNode} node
     * @param {string} name
     */
    function verifyName(node, name) {
      for (const option of options) {
        if (option.test(name)) {
          context.report({
            node,
            messageId: 'restrictedName',
            data: {
              message:
                option.message || `Using \`${name}\` name is not allowed.`
            }
          })
          break
        }
      }
    }

    /**
     * @param {Literal | TemplateLiteral} node
     */
    function verifyNameNode(node) {
      /** @type {string | undefined} */
      let name
      if (node.type === 'TemplateLiteral') {
        const quasis = node.quasis[0]
        name = quasis.value.cooked
      } else {
        name = `${node.value}`
      }
      if (name != null) {
        verifyName(node, name)
      }
    }

    return utils.compositingVisitors(
      utils.executeOnCallVueComponent(context, (node) => {
        if (node.arguments.length === 2) {
          const argument = node.arguments[0]

          if (canVerify(argument)) {
            verifyNameNode(argument)
          }
        }
      }),
      utils.executeOnVue(context, (obj) => {
        // Report if a component has been registered locally with a restricted name.
        for (const { node, name } of utils.getRegisteredComponents(obj)) {
          verifyName(node, name)
        }

        const node = utils.findProperty(obj, 'name')

        if (!node) return
        if (!canVerify(node.value)) return
        verifyNameNode(node.value)
      }),
      utils.defineScriptSetupVisitor(context, {
        onDefineOptionsEnter(node) {
          if (node.arguments.length === 0) return
          const define = node.arguments[0]
          if (define.type !== 'ObjectExpression') return
          const nameNode = utils.findProperty(define, 'name')
          if (!nameNode) return
          if (!canVerify(nameNode.value)) return
          verifyNameNode(nameNode.value)
        }
      })
    )
  }
}
