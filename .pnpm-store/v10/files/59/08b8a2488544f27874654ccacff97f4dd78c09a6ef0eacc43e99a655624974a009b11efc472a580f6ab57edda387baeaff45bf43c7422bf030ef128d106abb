type Factory<Key, Value> = (key: Key) => Value;
declare class CachedFactory<Key, Value> {
    #private;
    constructor(factory: Factory<Key, Value>);
    clear(): void;
    entries(): IterableIterator<[Key, Value]>;
    get(key: Key): Value;
}

export { CachedFactory, Factory };
