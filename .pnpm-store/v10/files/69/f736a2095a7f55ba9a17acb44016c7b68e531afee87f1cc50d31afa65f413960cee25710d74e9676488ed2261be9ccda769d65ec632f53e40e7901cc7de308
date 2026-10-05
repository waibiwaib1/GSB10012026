import * as eslint from 'eslint';
import * as jsonc_eslint_parser from 'jsonc-eslint-parser';
import { PackageJsonRuleContext } from '../createRule.js';
import 'estree';

type Options = [{
    order: Order;
}?];
type Order = "legacy" | "sort-package-json" | string[];
declare const rule: {
    create(context: PackageJsonRuleContext<Options>): jsonc_eslint_parser.RuleListener;
    meta: eslint.Rule.RuleMetaData;
};

export { rule };
