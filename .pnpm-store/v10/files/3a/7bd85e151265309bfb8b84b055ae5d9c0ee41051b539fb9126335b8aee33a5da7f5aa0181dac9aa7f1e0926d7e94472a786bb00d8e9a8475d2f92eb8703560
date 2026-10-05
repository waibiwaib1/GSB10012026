import { EmojiPhysics, EmojiPosition, EmojiProcess } from "./actor.js";
import { EmojiTick } from "./animate.js";
import { EmojiEvents } from "./events.js";
import { MakePartial } from "./utils.js";
/**
 * Settings to launch a blast of emojis! 🎆
 */
export interface EmojiBlastSettings {
    /**
     * Class name to add to all emoji elements.
     */
    className: string;
    /**
     * Element container to append elements into.
     */
    container: SettingValue<Element>;
    /**
     * How many emojis to create per blast.
     */
    emojiCount: SettingValue<number>;
    /**
     * Allowed potential emoji names to set as textContent.
     */
    emojis: SettingValue<string[]>;
    /**
     * Handlers for user interactions with individual emojis.
     */
    events: EmojiEvents;
    /**
     * Runtime change constants for emoji element movements.
     */
    physics: MakePartial<EmojiPhysics, "initialVelocities">;
    /**
     * How to determine where to place blasts of emojis around the page.
     */
    position: SettingValue<EmojiPosition>;
    /**
     * Processes each element just before it's appended to the container.
     */
    process: EmojiProcess;
    /**
     * Hook to call on each tick.
     */
    tick?: EmojiTick;
    /**
     * How many different types of emojis are allowed within a blast.
     */
    uniqueness: SettingValue<number>;
}
/**
 * Setting value or a method to create it.
 * @template T   Type of the setting value.
 */
export type SettingValue<T> = (() => T) | T;
/**
 * Default class name to add to emoji elements.
 */
export declare const defaultClassName = "emoji-styles";
/**
 * Default creator for a container element.
 * @returns <div /> element prepended to document.body.
 */
export declare const defaultCreateContainer: () => HTMLElement;
/**
 * Default emojiCount to choose a random number of emoji per blast.
 * @returns Random integer within 14 to 28.
 */
export declare const defaultEmojiCount: () => number;
export declare const defaultEvents: EmojiEvents;
/**
 * Default runtime change constants for actor movements.
 */
export declare const defaultPhysics: EmojiPhysics;
/**
 * Default position to choose random locations within the page.
 * @returns Random { left, top } integers within the page.
 */
export declare const defaultPosition: () => {
    x: number;
    y: number;
};
/**
 * Launches a blast of emojis across the page! 🎆
 */
export declare const emojiBlast: (settings?: Partial<EmojiBlastSettings>) => {
    stop: () => void;
};
