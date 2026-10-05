import { EmojiBlastSettings, SettingValue } from "./emojiBlast.js";
/**
 * Settings to periodically blast emojis! 🎆
 */
export interface EmojiBlastsSettings extends EmojiBlastSettings {
    /**
     * How frequently to create blasts.
     */
    interval: SettingValue<number>;
    /**
     * Schedules blasts to occur.
     */
    scheduler: EmojiScheduler;
}
/**
 * Schedules an blast to occur.
 * @param action   Action that causes the blast.
 * @param delay   How long before the action should occur.
 */
export type EmojiScheduler = (action: () => void, delay: number) => void;
/**
 * Returned handler for an ongoing blasts of emojis run.
 */
export interface EmojiBlastsHandler {
    /**
     * Triggers a blast of emojis.
     */
    blast: () => void;
    /**
     * Cancels the emojiBlasts run.
     */
    cancel: () => void;
}
/**
 * Periodically blast emojis across the page! 🎆
 * @param settings   Settings to blast emojis.
 * @returns Handler for the ongoing blasts of emojis.
 */
export declare const emojiBlasts: (settings?: Partial<EmojiBlastsSettings>) => EmojiBlastsHandler;
