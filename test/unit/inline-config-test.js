'use strict';

var assert = require('power-assert');
var Linter = require('../../lib');

function linterFor(ruleName) {
  var rules = {};
  rules[ruleName] = true;

  return new Linter({ config: { rules: rules } });
}

function violations(ruleName, template) {
  var linter = linterFor(ruleName);

  return linter.verify({ source: template, moduleId: 'layout.hbs' });
}

describe('inline template-lint configuration', function() {
  describe('child-less mustache comments', function() {
    it('disables all rules for later siblings and their descendants', function() {
      var template = '<div>\n  {{! template-lint disabled=true }}\n  {{{foo}}}\n  <span>{{{bar}}}</span>\n</div>';

      assert.deepEqual(violations('triple-curlies', template), []);
    });

    it('still lints earlier siblings', function() {
      var template = '<div>\n  {{{before}}}\n  {{! template-lint disabled=true }}\n  {{{foo}}}\n</div>';

      assert.equal(violations('triple-curlies', template).length, 1);
    });

    it('does not apply to siblings after the containing block closes', function() {
      var template = '<div>\n  {{! template-lint disabled=true }}\n  {{{foo}}}\n</div>\n{{{baz}}}';

      assert.equal(violations('triple-curlies', template).length, 1);
    });

    it('disables a single rule by name within the block', function() {
      var template = '{{#foo-bar as |baz|}}\n  {{! template-lint triple-curlies=false }}\n  {{{inside}}}\n{{/foo-bar}}\n{{{outside}}}';

      assert.equal(violations('triple-curlies', template).length, 1);
    });

    it('does not disable other rules when a rule name is given', function() {
      var template = '{{#foo-bar as |baz|}}\n  {{! template-lint triple-curlies=false }}\n  plain inside\n{{/foo-bar}}';

      assert.equal(violations('bare-strings', template).length, 1);
    });
  });

  describe('in-element mustache comments', function() {
    it('does not disable a rule for descendants unless recursive=true', function() {
      var template = '<div {{! template-lint rule=\'triple-curlies\' disabled=true}}>{{{foo}}}</div>';

      assert.equal(violations('triple-curlies', template).length, 1);
    });

    it('still disables a rule for the element itself without recursive=true', function() {
      var template = '<div style="foo{{bar}}" {{! template-lint rule=\'style-concatenation\' disabled=true}}></div>';

      assert.deepEqual(violations('style-concatenation', template), []);
    });

    it('disables a rule for descendants when recursive=true', function() {
      var template = '<div {{! template-lint rule=\'triple-curlies\' disabled=true recursive=true}}>{{{foo}}}</div>';

      assert.deepEqual(violations('triple-curlies', template), []);
    });

    it('only disables the rule named via rule=', function() {
      var template = '<div {{! template-lint rule=\'bare-strings\' disabled=true recursive=true}}>{{{foo}}}</div>';

      assert.equal(violations('triple-curlies', template).length, 1);
    });

    it('does not leak outside of the element when recursive=true', function() {
      var template = '<div {{! template-lint rule=\'triple-curlies\' disabled=true recursive=true}}>{{{foo}}}</div>\n{{{bar}}}';

      assert.equal(violations('triple-curlies', template).length, 1);
    });
  });

  describe('legacy html comments', function() {
    it('disables a rule for the rest of the template', function() {
      var template = '<!-- template-lint triple-curlies=false -->\n{{{foo}}}';

      assert.deepEqual(violations('triple-curlies', template), []);
    });

    it('supports enabled=false for all rules', function() {
      var template = '<!-- template-lint enabled=false -->\n{{{foo}}}';

      assert.deepEqual(violations('triple-curlies', template), []);
    });
  });
});
