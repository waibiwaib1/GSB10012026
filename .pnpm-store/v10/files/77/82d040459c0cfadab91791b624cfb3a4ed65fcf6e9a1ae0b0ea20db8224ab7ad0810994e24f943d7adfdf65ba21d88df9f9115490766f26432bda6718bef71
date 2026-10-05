import { emojiBlasts } from "emoji-blast";
// https://github.com/JoshuaKGoldberg/emoji-blast/issues/822#issuecomment-2532742672
// @ts-expect-error -- Working around CJS/ESM default interop
import * as KonamiCodeJS from "konami-code-js";
// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
const KonamiCode = 
// eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
KonamiCodeJS.default ?? KonamiCodeJS;
/**
 * @returns Function that stops emojis when called.
 */
export const initializeKonamiEmojiBlast = (onActivate) => {
    let explosions;
    const konamiCode = new KonamiCode(() => {
        if (explosions) {
            explosions.cancel();
            explosions = undefined;
        }
        else {
            explosions = emojiBlasts();
            onActivate?.();
        }
    });
    return () => {
        explosions?.cancel();
        konamiCode.disable();
    };
};
//# sourceMappingURL=initializeKonamiEmojiBlast.js.map