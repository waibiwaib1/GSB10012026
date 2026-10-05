import { AST } from 'jsonc-eslint-parser';

declare function isJSONStringLiteral(node: AST.JSONNode): node is AST.JSONStringLiteral;
declare function isNotNullish<T extends NonNullable<unknown>>(value: null | T | undefined): value is T;

export { isJSONStringLiteral, isNotNullish };
