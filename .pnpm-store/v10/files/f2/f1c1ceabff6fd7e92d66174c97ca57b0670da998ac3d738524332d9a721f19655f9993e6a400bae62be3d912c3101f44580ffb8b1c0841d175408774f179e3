import * as ESTree from 'estree';
import { AST as AST$1, Rule, SourceCode } from 'eslint';
import { AST, RuleListener } from 'jsonc-eslint-parser';

type JsonAstBodyExpression = ESTree.Expression & {
    properties: JsonAstBodyProperty[];
};
type JsonAstBodyProperty = AST.JSONProperty & {
    value: string;
};
interface JsonAstBodyStatement extends ESTree.ExpressionStatement {
    expression: JsonAstBodyExpression;
}
interface PackageJsonAst extends AST$1.Program {
    body: [JsonAstBodyStatement];
}
interface PackageJsonRuleContext<Options extends unknown[] = unknown[]> extends Rule.RuleContext {
    options: Options;
    sourceCode: PackageJsonSourceCode;
}
interface PackageJsonRuleModule<Options extends unknown[] = unknown[]> {
    create(context: PackageJsonRuleContext<Options>): RuleListener;
    meta: Rule.RuleMetaData;
}
interface PackageJsonSourceCode extends SourceCode {
    ast: PackageJsonAst;
}
declare function createRule<Options extends unknown[]>(rule: PackageJsonRuleModule<Options>): {
    create(context: PackageJsonRuleContext<Options>): RuleListener;
    meta: Rule.RuleMetaData;
};

export { type JsonAstBodyExpression, type JsonAstBodyProperty, type JsonAstBodyStatement, type PackageJsonAst, type PackageJsonRuleContext, type PackageJsonRuleModule, type PackageJsonSourceCode, createRule };
