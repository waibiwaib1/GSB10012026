/**
 * Like Partial, but also makes particular properties Partial one level deep.
 */
export type MakePartial<T, PartialKeys extends keyof T> = {
    [K in PartialKeys]?: Partial<T[K]>;
} & Partial<Omit<T, PartialKeys>>;
/**
 * Grabs the value of an item or item-returning function.
 * @param value   Item or item-returning function.
 */
export declare const obtainValue: <T>(value: (() => T) | T) => T;
/**
 * Grabs a single random member of an array.
 * @template T   Type of items in the array.
 * @param array   Array of items.
 */
export declare const randomArrayMember: <T>(array: readonly T[]) => T;
/**
 * Creates a shuffled version of an array.
 * @template T   Type of items in the array.
 * @param array   Array to copy.
 * @returns Shuffled version of the array.
 */
export declare const shuffleArray: <T>(array: readonly T[]) => T[];
