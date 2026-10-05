import { emojiBlast } from "./emojiBlast.js";
import { obtainValue } from "./utils.js";
/**
 * Default interval setting for fire delays.
 * @returns Random number between 0 and 2100.
 */
const defaultInterval = () => 700 + Math.floor(Math.random() * 1401);
/**
 * Periodically blast emojis across the page! 🎆
 * @param settings   Settings to blast emojis.
 * @returns Handler for the ongoing blasts of emojis.
 */
export const emojiBlasts = (settings = {}) => {
    const { interval = defaultInterval, scheduler = setTimeout } = settings;
    let cancelled = false;
    const blast = () => {
        emojiBlast(settings);
    };
    const blastAndSchedule = () => {
        if (cancelled) {
            return;
        }
        if (document.visibilityState === "visible") {
            blast();
        }
        scheduler(blastAndSchedule, obtainValue(interval));
    };
    scheduler(blastAndSchedule, 0);
    return {
        blast,
        cancel() {
            cancelled = true;
        },
    };
};
//# sourceMappingURL=emojiBlasts.js.map