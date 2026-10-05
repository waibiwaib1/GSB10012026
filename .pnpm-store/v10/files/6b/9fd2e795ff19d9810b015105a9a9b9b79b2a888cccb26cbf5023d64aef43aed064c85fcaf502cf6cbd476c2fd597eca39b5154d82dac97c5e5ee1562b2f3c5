import { EmojiActor } from "./actor.js";
/**
 * Returned handler for an ongoing blasts of emojis run.
 */
export interface EmojiBlastHandler {
    /**
     * Stops physics in the emojiBlast ticks.
     */
    stop: () => void;
}
/**
 * Hook to call on each tick.
 * @param actors   Each actor in play at the time.
 */
export type EmojiTick = (actors: EmojiActor[]) => void;
/**
 * Starts the regular gameplay loop of telling actors to animate.
 *
 * Each game "tick" is scheduled using `requestAnimationFrame`.
 * During each tick, each actor is told to `act` with the time elapsed.
 * If it indicates that it's out of bounds, it's removed from the actors array.
 */
export declare function animate(
/**
 * Actors that have been added and not yet marked themselves as out of bounds.
 */
actors: EmojiActor[], 
/**
 * Handler to run before each tick, if provided
 */
beforeTick: EmojiTick | undefined): {
    stop: () => void;
};
