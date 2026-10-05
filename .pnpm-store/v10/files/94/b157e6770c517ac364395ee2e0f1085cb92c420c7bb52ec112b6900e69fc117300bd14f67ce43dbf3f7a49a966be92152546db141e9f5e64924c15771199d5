import { randomInRange } from "./range.js";
import { randomArrayMember } from "./utils.js";
/**
 * Pixel distance out of the screen bounds to treat actors as out-of-bounds.
 */
const outOfBounds = 350;
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
export class EmojiActor {
    /**
     * CSS opacity style, starting at 1 for fully visible.
     */
    #opacity = 1;
    /**
     * Runtime change constants for actor movements.
     */
    #physics;
    /**
     * Current element coordinates and rotation.
     */
    #position;
    /**
     * Change amounts for element position.
     */
    #velocity;
    /**
     * Attached element kept in the DOM.
     */
    element;
    constructor(settings) {
        this.element = document.createElement("button");
        this.element.className = settings.className;
        this.element.textContent = randomArrayMember(settings.emojis);
        this.element.setAttribute("role", "img");
        this.element.style.fontSize = `${randomInRange(settings.physics.fontSize)}px`;
        this.element.style.transition = "16ms opacity, 16ms transform";
        this.#physics = settings.physics;
        this.#position = {
            rotation: randomInRange(settings.physics.rotation),
            x: settings.position.x,
            y: settings.position.y,
        };
        this.#velocity = {
            rotation: randomInRange(settings.physics.initialVelocities.rotation),
            x: randomInRange(settings.physics.initialVelocities.x),
            y: randomInRange(settings.physics.initialVelocities.y),
        };
        this.updateElement();
        settings.process?.(this.element);
        settings.container.appendChild(this.element);
    }
    /**
     * Updates the attached DOM element to match tracking position.
     */
    updateElement() {
        this.element.style.opacity = `${this.#opacity}`;
        this.element.style.transform = `translate(${this.#position.x}px, ${this.#position.y}px) rotate(${Math.round(this.#position.rotation)}deg)`;
    }
    /**
     * Moves the actor forward one tick.
     * @param timeElapsed   How many milliseconds have passed since the last action.
     * @returns Whether this is now dead.
     */
    act(timeElapsed) {
        if (this.#physics.opacityDecay) {
            this.#opacity -=
                timeElapsed / (this.#physics.opacityDecay * this.#physics.framerate);
            if (this.#opacity <= 0) {
                return true;
            }
        }
        this.#velocity.rotation *= this.#physics.rotationDeceleration;
        this.#velocity.y += this.#physics.gravity;
        this.#position.rotation += this.#velocity.rotation;
        this.#position.x +=
            (this.#velocity.x * timeElapsed) / this.#physics.framerate;
        this.#position.y +=
            (this.#velocity.y * timeElapsed) / this.#physics.framerate;
        const windowHeight = window.outerHeight || document.documentElement.clientHeight;
        const windowWidth = window.outerWidth || document.documentElement.clientWidth;
        if (!this.#physics.preserveOutOfBounds) {
            if (this.#position.y - this.element.clientHeight >
                windowHeight + outOfBounds) {
                return true;
            }
            if (this.#position.y + this.element.clientHeight < -outOfBounds) {
                return true;
            }
            if (this.#position.x - this.element.clientWidth >
                windowWidth + outOfBounds) {
                return true;
            }
            if (this.#position.x + this.element.clientWidth < -outOfBounds) {
                return true;
            }
        }
        this.updateElement();
        return false;
    }
    /**
     * Disposes of the attached DOM element upon actor death.
     */
    dispose() {
        if (this.element.parentElement !== null) {
            this.element.parentElement.removeChild(this.element);
        }
    }
    /**
     * Updates the emoji for being clicked.
     */
    update(updates) {
        if (updates.opacity !== undefined) {
            this.#opacity = updates.opacity;
        }
        if (updates.velocity !== undefined) {
            if (updates.velocity.rotation !== undefined) {
                this.#velocity.rotation = updates.velocity.rotation;
            }
            if (updates.velocity.x !== undefined) {
                this.#velocity.x = updates.velocity.x;
            }
            if (updates.velocity.y !== undefined) {
                this.#velocity.y = updates.velocity.y;
            }
        }
    }
    /**
     * CSS opacity style, starting at 1 for fully visible.
     */
    get opacity() {
        return this.#opacity;
    }
    /**
     * Current element coordinates and rotation.
     */
    get position() {
        return this.#position;
    }
    /**
     * Change amounts for element position.
     */
    get velocity() {
        return this.#velocity;
    }
}
//# sourceMappingURL=actor.js.map