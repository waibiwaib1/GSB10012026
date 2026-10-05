/**
 * DOM attribute indicating that events were initialized for a container.
 */
const attributeIndicator = "data-emoji-blast-events-initialized";
const domNodesToActors = new WeakMap();
export function initializeEvents(actors, container, events) {
    for (const actor of actors) {
        domNodesToActors.set(actor.element, actor);
    }
    if (container.hasAttribute(attributeIndicator)) {
        return;
    }
    container.setAttribute(attributeIndicator, "true");
    container.addEventListener("click", (event) => {
        const actor = event.target && domNodesToActors.get(event.target);
        if (actor) {
            events.onClick({ actor, event });
        }
    });
}
//# sourceMappingURL=events.js.map