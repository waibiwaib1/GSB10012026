import { AST } from 'jsonc-eslint-parser';

type JSONPropertyWithKeyAndValue<Value extends string> = AST.JSONProperty & {
    key: AST.JSONStringLiteral;
    value: Value;
};
declare function findPropertyWithKeyValue<Value extends string>(properties: AST.JSONProperty[], value: Value): JSONPropertyWithKeyAndValue<Value> | undefined;

export { type JSONPropertyWithKeyAndValue, findPropertyWithKeyValue };
