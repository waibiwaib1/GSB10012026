import * as eslint from 'eslint';
import * as jsonc_eslint_parser from 'jsonc-eslint-parser';
import { PackageJsonRuleContext } from '../createRule.js';
import 'estree';

/**
 * Given a top-level property name, create a rule that requires that property to be present.
 * Optionally, include it in the recommended config.
 */
declare const createRequirePropertyRule: (propertyName: string, isRecommended?: boolean) => {
    create(context: PackageJsonRuleContext<unknown[]>): jsonc_eslint_parser.RuleListener;
    meta: eslint.Rule.RuleMetaData;
};

export { createRequirePropertyRule };
