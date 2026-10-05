import { RandomRange } from "./range.js";
/**
 * Settings to create a single emoji within a container.
 */
export interface EmojiActorSettings {
    /**
     * Class name to add to the actor's element.
     */
    className: string;
    /**
     * Element container to append the element into.
     */
    container: Element;
    /**
     * Allowed potential emoji names to set as textContent.
     */
    emojis: string[];
    /**
     * Runtime change constants for actor movements.
     */
    physics: EmojiPhysics;
    /**
     * How to determine where to place blasts of emojis around the page.
     */
    position: EmojiPosition;
    /**
     * Processes each element just before it's appended to the container.
     */
    process?: EmojiProcess;
}
/**
 * Runtime change constants for actor movements.
 */
export interface EmojiPhysics {
    /**
     * Individual emojis' font size range.
     */
    fontSize: RandomRange;
    /**
     * Expected frames per second to adjust position and velocity changes by.
     */
    framerate: number;
    /**
     * How much to increase y-velocity downward each tick.
     */
    gravity: number;
    /**
     * Initial velocity ranges for individual emojis.
     */
    initialVelocities: InitialVelocities;
    /**
     * How much to slow down the (time elapsed / framerate) opacity reduction each tick.
     */
    opacityDecay?: number;
    /**
     * Whether to skip removing emojis that move outside of the visible screen.
     */
    preserveOutOfBounds?: boolean;
    /**
     * Individual emojis' initial rotation range.
     */
    rotation: RandomRange;
    /**
     * How much to decrease rotation amount each tick.
     */
    rotationDeceleration: number;
}
/**
 * Initial velocity ranges for individual emojis.
 */
export interface InitialVelocities {
    /**
     * Range of initial rotation amount.
     */
    rotation: RandomRange;
    /**
     * Range of initial horizontal velocity.
     */
    x: RandomRange;
    /**
     * Range of initial vertical velocity.
     */
    y: RandomRange;
}
/**
 * Absolute CSS position to place an emoji element at.
 */
export interface EmojiPosition {
    /**
     * Pixels to offset by the left.
     */
    x: number;
    /**
     * Pixels to offset by the top.
     */
    y: number;
}
/**
 * In-progress tracking for an actor's position.
 */
export type EmojiVelocity = {
    /**
     * How much the actor's element is rotated.
     */
    rotation: number;
} & EmojiPosition;
/**
 * Processes an element just before it's appended to the container.
 * @param element   Element about to be appended to the container.
 */
export type EmojiProcess = (element: Element) => void;
/**
 * Updates to apply to an emoji's position.
 */
export interface EmojiUpdates {
    /**
     * Updated emoji opacity, if it should change.
     */
    opacity?: number;
    /**
     * Any velocity changes, if they should change.
     */
    velocity?: Partial<EmojiVelocity>;
}
/**
 * Contains the position, velocity, and DOM element for a single emoji.
 *
 * This creates and keeps a single DOM element button in the DOM.
 * Text content for the button is determined by the provided actors.
 *
 * On each game tick, this actor will:
 *  1. Dispose itself if it's moved past out of the game screen
 *  2. Reduce opacity a little bit
 *  3. Dispose itself if it's no longer visible at all
 *  4. Adjust position and velocity as per its physics constants
 *  5. Update the DOM element's opacity and position to reflect those changes
 *
 * "Disposing" an actor means removing its element from the document.
 */
export declare class EmojiActor {
    #private;
    /**
     * Attached element kept in the DOM.
     */
    readonly element: HTMLSpanElement;
    constructor(settings: EmojiActorSettings);
    /**
     * Updates the attached DOM element to match tracking position.
     */
    private updateElement;
    /**
     * Moves the actor forward one tick.
     * @param timeElapsed   How many milliseconds have passed since the last action.
     * @returns Whether this is now dead.
     */
    act(timeElapsed: number): boolean;
    /**
     * Disposes of the attached DOM element upon actor death.
     */
    dispose(): void;
    /**
     * Updates the emoji for being clicked.
     */
    update(updates: EmojiUpdates): void;
    /**
     * CSS opacity style, starting at 1 for fully visible.
     */
    get opacity(): number;
    /**
     * Current element coordinates and rotation.
     */
    get position(): Readonly<EmojiPosition>;
    /**
     * Change amounts for element position.
     */
    get velocity(): Readonly<EmojiVelocity>;
}
