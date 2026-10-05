import { EmojiActor, } from "./actor.js";
import { animate } from "./animate.js";
import { defaultEmojis } from "./emojis.js";
import { initializeEvents } from "./events.js";
import { createStyleElementAndClass } from "./styles.js";
import { obtainValue, shuffleArray } from "./utils.js";
/**
 * Default class name to add to emoji elements.
 */
export const defaultClassName = "emoji-styles";
/**
 * Default creator for a container element.
 * @returns <div /> element prepended to document.body.
 */
export const defaultCreateContainer = (() => {
    let container;
    return () => {
        if (container?.parentNode === document.body) {
            return container;
        }
        container = document.createElement("div");
        document.body.prepend(container);
        return container;
    };
})();
/**
 * Default emojiCount to choose a random number of emoji per blast.
 * @returns Random integer within 14 to 28.
 */
export const defaultEmojiCount = () => Math.floor(Math.random() * 14) + 14;
export const defaultEvents = {
    onClick({ actor }) {
        actor.update({
            opacity: 1,
            velocity: {
                y: actor.velocity.y / 2 - 15,
            },
        });
    },
};
/**
 * Default runtime change constants for actor movements.
 */
export const defaultPhysics = {
    fontSize: {
        max: 28,
        min: 14,
    },
    framerate: 60,
    gravity: 0.35,
    initialVelocities: {
        rotation: {
            max: 7,
            min: -7,
        },
        x: {
            max: 7,
            min: -7,
        },
        y: {
            max: -7,
            min: -21,
        },
    },
    preserveOutOfBounds: false,
    rotation: {
        max: 45,
        min: -45,
    },
    rotationDeceleration: 0.98,
};
/**
 * Default position to choose random locations within the page.
 * @returns Random { left, top } integers within the page.
 */
export const defaultPosition = () => ({
    x: Math.random() * innerWidth,
    y: Math.random() * innerHeight,
});
/**
 * Launches a blast of emojis across the page! 🎆
 */
export const emojiBlast = (settings = {}) => {
    const { className = defaultClassName, container: containerSetting = defaultCreateContainer, emojiCount = defaultEmojiCount, emojis = defaultEmojis, events = defaultEvents, position = defaultPosition, process, tick, uniqueness = Infinity, } = settings;
    const container = obtainValue(containerSetting);
    createStyleElementAndClass(className);
    const physics = {
        ...defaultPhysics,
        ...settings.physics,
        initialVelocities: {
            ...defaultPhysics.initialVelocities,
            ...(settings.physics !== undefined
                ? settings.physics.initialVelocities
                : {}),
        },
    };
    const emojiSettings = {
        className,
        container,
        // Copy the input array to prevent modifications.
        emojis: shuffleArray(obtainValue(emojis)).slice(0, obtainValue(uniqueness)),
        physics,
        position: obtainValue(position),
        process,
    };
    const blastEmojiCount = obtainValue(emojiCount);
    const actors = [];
    for (let i = 0; i < blastEmojiCount; i += 1) {
        actors.push(new EmojiActor(emojiSettings));
    }
    initializeEvents(actors, container, events);
    return animate(actors, tick);
};
//# sourceMappingURL=emojiBlast.js.map