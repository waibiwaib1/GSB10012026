function withoutUndefinedPropertiesDeep(value) {
  const result = { ...value };
  for (const k in result) {
    if (result[k] === void 0) {
      delete result[k];
    } else if (result[k] !== null && typeof result[k] === "object") {
      result[k] = withoutUndefinedPropertiesDeep(result[k]);
    }
  }
  return result;
}
export {
  withoutUndefinedPropertiesDeep
};
//# sourceMappingURL=withoutUndefinedPropertiesDeep.js.map