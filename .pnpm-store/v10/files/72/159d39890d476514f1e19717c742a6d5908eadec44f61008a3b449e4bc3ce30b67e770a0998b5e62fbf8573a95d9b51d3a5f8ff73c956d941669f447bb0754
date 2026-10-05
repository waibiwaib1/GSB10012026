class CachedFactory {
  #cache = /* @__PURE__ */ new Map();
  #getter;
  constructor(factory) {
    this.#getter = factory;
  }
  clear() {
    this.#cache.clear();
  }
  entries() {
    return this.#cache.entries();
  }
  get(key) {
    const existing = this.#cache.get(key);
    if (existing) {
      return existing;
    }
    const value = this.#getter(key);
    this.#cache.set(key, value);
    return value;
  }
}
export {
  CachedFactory
};
//# sourceMappingURL=index.js.map