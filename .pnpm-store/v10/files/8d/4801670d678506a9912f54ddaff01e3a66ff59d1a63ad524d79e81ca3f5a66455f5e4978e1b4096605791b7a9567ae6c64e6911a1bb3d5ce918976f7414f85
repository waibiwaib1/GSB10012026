import { RuleTester } from 'eslint';
import { PackageJsonRuleModule } from '../../createRule.js';
import 'estree';
import 'jsonc-eslint-parser';

type JsonRuleTester = RuleTester & {
    run: JsonRuleTesterRun;
};
type JsonRuleTesterRun = (name: string, rule: PackageJsonRuleModule, tests: {
    invalid?: RuleTester.InvalidTestCase[] | undefined;
    valid?: (RuleTester.ValidTestCase | string)[] | undefined;
}) => void;
declare const ruleTester: JsonRuleTester;

export { type JsonRuleTester, type JsonRuleTesterRun, ruleTester };
