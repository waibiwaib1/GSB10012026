'use strict';

var assign = require('lodash').assign;

var TEMPLATE_LINT_PREFIX = 'template-lint';

function parseConfigDirective(rawValue) {
  var value = rawValue.trim();

  if (value.indexOf(TEMPLATE_LINT_PREFIX) !== 0) {
    return null;
  }

  var hash = {};
  var hashPartPattern = /([\w-]+)\s*=\s*(?:"([^"]*)"|'([^']*)'|(\S+))/g;
  var match;

  while ((match = hashPartPattern.exec(value)) !== null) {
    hash[match[1]] = match[2] || match[3] || match[4];
  }

  return hash;
}

function invokeVisitorEnter(visitor, args, thisContext) {
  var func;

  if (typeof visitor === 'function') {
    func = visitor;
  } else if (visitor && visitor.enter) {
    func = visitor.enter;
  }

  return func.apply(thisContext, args);
}

module.exports = function(options) {
  var log = options.log;
  var config = options.config;
  var ruleName = options.name;
  var defaultSeverity = options.defaultSeverity;

  function BasePlugin(options) {
    this.options = options;
    this.syntax = null; // set by Glimmer

    // split into a source array (allow windows and posix line endings)
    this.source = this.options.rawSource.split(/(?:\r\n?|\n)/g);

    this._log = log;
    this.ruleName = ruleName;
    this.severity = defaultSeverity;
    this.config = this.parseConfig(config);
    this._configStack = [];
  }

  BasePlugin.prototype.parseConfig = function(config) {
    return config;
  };

  BasePlugin.prototype.transform = function(ast) {
    this.syntax.traverse(ast, this.getVisitors());

    return ast;
  };


  BasePlugin.prototype.getVisitors = function() {
    var pluginContext = this;
    var visitors = {};
    var ruleVisitors = this.visitors();
    var ruleVisitor;

    for (var key in ruleVisitors) {
      ruleVisitor = ruleVisitors[key];

      visitors[key] = this.processVisitor(ruleVisitor);
    }

    ['Program', 'ElementNode', 'BlockStatement'].forEach(function(nodeType) {
      var existingVisitor = visitors[nodeType];
      var existingEnter = existingVisitor && existingVisitor.enter;
      var existingExit = existingVisitor && existingVisitor.exit;

      visitors[nodeType] = {
        enter: function(node) {
          pluginContext._pushConfigFrame(node);

          if (existingEnter) {
            existingEnter.apply(pluginContext, arguments);
          }
        },

        exit: function() {
          if (existingExit) {
            existingExit.apply(pluginContext, arguments);
          }

          pluginContext._popConfigFrame();
        }
      };
    });

    visitors['MustacheCommentStatement'] = this.processVisitor(function(node) {
      pluginContext._processInlineConfigNode(node);
    });

    var commentVisitor = visitors['CommentStatement'];
    if (!commentVisitor) {
      commentVisitor = { enter: function() {} };
    } else if (typeof commentVisitor === 'function') {
      commentVisitor = { enter: commentVisitor };
    }

    var commentVisitorEnter = commentVisitor.enter;
    commentVisitor.enter = function(node) {
      pluginContext._processLegacyConfigNode(node);

      if (commentVisitorEnter) {
        commentVisitorEnter.apply(pluginContext, arguments);
      }
    };
    visitors['CommentStatement'] = commentVisitor;

    return visitors;
  };

  BasePlugin.prototype.processVisitor = function(ruleVisitor) {
    var pluginContext = this;
    var visitor = {
      enter: function(node) {
        if (!pluginContext._isNodeDisabled(node)) {
          invokeVisitorEnter(ruleVisitor, [node], pluginContext);
        }
      }
    };

    if (ruleVisitor && ruleVisitor.exit) {
      visitor.exit = function(node) {
        if (!pluginContext._isNodeDisabled(node)) {
          ruleVisitor.exit.apply(pluginContext, [node]);
        }
      };
    }

    return visitor;
  };

  BasePlugin.prototype.visitors = function() { };

  BasePlugin.prototype._currentConfig = function() {
    if (this._configStack.length === 0) {
      return this.config;
    }

    return this._configStack[this._configStack.length - 1].config;
  };

  BasePlugin.prototype._isRuleDisabled = function(hash) {
    if (hash.rule !== undefined) {
      return hash.rule === ruleName && hash.disabled === 'true';
    }

    return hash.disabled === 'true' ||
      hash.enable === 'false' ||
      hash.disable === 'true' ||
      hash.enabled === 'false' ||
      hash[ruleName] === 'false';
  };

  BasePlugin.prototype._pushConfigFrame = function(node) {
    var parentFrame = this._configStack[this._configStack.length - 1];
    var parentConfig = parentFrame ? parentFrame.config : this.config;
    var frame = {
      config: parentConfig,
      siblingDisabled: !!(parentFrame && (parentFrame.siblingDisabled || parentFrame.recursiveElementDisabled)),
      elementDisabled: false,
      recursiveElementDisabled: false,
      nodes: []
    };

    this._configStack.push(frame);

    if (node.type === 'ElementNode') {
      this._processElementConfigComments(node, frame);
    }
  };

  BasePlugin.prototype._popConfigFrame = function() {
    this._configStack.pop();
  };

  BasePlugin.prototype._processElementConfigComments = function(elementNode, frame) {
    if (!elementNode.comments) {
      return;
    }

    for (var i = 0; i < elementNode.comments.length; i++) {
      var hash = parseConfigDirective(elementNode.comments[i].value);

      if (hash && this._isRuleDisabled(hash)) {
        frame.elementDisabled = true;

        if (hash.recursive === 'true') {
          frame.siblingDisabled = true;
          frame.recursiveElementDisabled = true;
        }
      }
    }
  };

  BasePlugin.prototype._processInlineConfigNode = function(node) {
    var frame = this._configStack[this._configStack.length - 1];

    if (!frame) {
      return;
    }

    var hash = parseConfigDirective(node.value);

    if (hash && this._isRuleDisabled(hash)) {
      frame.config = false;
      frame.siblingDisabled = true;
    }
  };

  // HTML comments (<!-- template-lint ... -->) apply for the rest of the
  // template to preserve backwards compatibility.
  BasePlugin.prototype._processLegacyConfigNode = function(node) {
    var hash = parseConfigDirective(node.value);

    if (!hash) {
      return;
    }

    if (this._isRuleDisabled(hash)) {
      this.config = false;

      var frame = this._configStack[this._configStack.length - 1];
      if (frame) {
        frame.config = false;
        frame.siblingDisabled = true;
      }
    }
  };

  BasePlugin.prototype._isNodeDisabled = function(node) {
    if (!this.config) {
      return true;
    }

    var frameCount = this._configStack.length;

    if (frameCount === 0) {
      return false;
    }

    var frame = this._configStack[frameCount - 1];

    if (!frame.config) {
      return true;
    }

    if (node.type === 'AttrNode' ||
        node.type === 'ElementModifierStatement') {
      return frame.siblingDisabled || frame.elementDisabled;
    }

    if (node.type === 'MustacheCommentStatement' ||
        node.type === 'CommentStatement' ||
        node.type === 'TextNode') {
      return frame.siblingDisabled;
    }

    if (frame.siblingDisabled) {
      return true;
    }

    if (node.type === 'ElementNode') {
      return frame.elementDisabled;
    }

    return false;
  };

  BasePlugin.prototype.isDisabled = function() {
    return !this._currentConfig();
  };

  BasePlugin.prototype.log = function(result) {
    var defaults = {
      moduleId: this.options.moduleName,
      rule: this.ruleName,
      severity: this.severity
    };
    var reportedResult = assign({}, defaults, result);

    this._log(reportedResult);
  };

  BasePlugin.prototype.detect = function() {
    throw new Error('Must implemented #detect');
  };

  BasePlugin.prototype.process = function() {
    throw new Error('Must implemented #process');
  };

  // mostly copy/pasta from tildeio/htmlbars with a few tweaks:
  // https://github.com/tildeio/htmlbars/blob/v0.4.17/packages/htmlbars-syntax/lib/parser.js#L59-L90
  BasePlugin.prototype.sourceForNode = function(node) {
    if (!node.loc) { return; }

    var firstLine = node.loc.start.line - 1;
    var lastLine = node.loc.end.line - 1;
    var currentLine = firstLine - 1;
    var firstColumn = node.loc.start.column;
    var lastColumn = node.loc.end.column;
    var string = [];
    var line;

    while (currentLine < lastLine) {
      currentLine++;
      line = this.source[currentLine];

      if (currentLine === firstLine) {
        if (firstLine === lastLine) {
          string.push(line.slice(firstColumn, lastColumn));
        } else {
          string.push(line.slice(firstColumn));
        }
      } else if (currentLine === lastLine) {
        string.push(line.slice(0, lastColumn));
      } else {
        string.push(line);
      }
    }

    return string.join('\n');
  };

  return BasePlugin;
};
