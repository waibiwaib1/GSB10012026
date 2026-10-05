import { EmojiActor } from "./actor.js";
export interface EmojiEventData {
    /**
     * Actor being interacted with.
     */
    actor: EmojiActor;
    /**
     * Original triggering DOM event.
     */
    event: Event;
}
/**
 * Handler for a user interaction with an individual emoji.
 * @param actor
 */
export type EmojiEventHandler = (data: EmojiEventData) => void;
/**
 * Handlers for user interactions with individual emojis.
 */
export interface EmojiEvents {
    /**
     * Handler for a user clicking an emoji.
     */
    onClick: EmojiEventHandler;
}
export declare function initializeEvents(actors: EmojiActor[], container: Element, events: EmojiEvents): void;
