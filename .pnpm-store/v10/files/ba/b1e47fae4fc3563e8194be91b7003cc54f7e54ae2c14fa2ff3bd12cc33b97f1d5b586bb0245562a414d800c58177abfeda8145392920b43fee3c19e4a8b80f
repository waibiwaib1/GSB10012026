type WithoutUndefinedProperties<T> = {
    [K in keyof T]: Exclude<T[K], undefined>;
} & {};
type WithoutUndefinedPropertiesDeep<T> = {
    [K in keyof T]?: T[K] extends object ? WithoutUndefinedPropertiesDeep<T[K]> : Exclude<T[K], undefined>;
} & {};

export type { WithoutUndefinedProperties, WithoutUndefinedPropertiesDeep };
