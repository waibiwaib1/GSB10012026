import { describe, expect, it, vi } from "vitest";
import { initializeKonamiEmojiBlast } from "./initializeKonamiEmojiBlast.js";
let mockKonamiCode;
vi.mock("konami-code-js", () => ({
    default: class MockKonamiCode {
        mockLaunch;
        disable = vi.fn();
        constructor(mockLaunch) {
            this.mockLaunch = mockLaunch;
            // eslint-disable-next-line @typescript-eslint/no-this-alias
            mockKonamiCode = this;
        }
    },
}));
const mockEmojiBlastsCancel = vi.fn();
const mockEmojiBlastsStart = vi.fn(() => ({
    cancel: mockEmojiBlastsCancel,
}));
vi.mock("emoji-blast", () => ({
    get emojiBlasts() {
        return mockEmojiBlastsStart;
    },
}));
describe("initializeKonamiEmojiBlast", () => {
    it("starts emoji blasts and calls onActivate when the konami code is triggered a first time", () => {
        const onActivate = vi.fn();
        initializeKonamiEmojiBlast(onActivate);
        mockKonamiCode.mockLaunch();
        expect(mockEmojiBlastsStart).toHaveBeenCalled();
        expect(onActivate).toHaveBeenCalledTimes(1);
    });
    it("stops emoji blasts when the konami code is triggered a second time", () => {
        const onActivate = vi.fn();
        initializeKonamiEmojiBlast(onActivate);
        mockKonamiCode.mockLaunch();
        mockKonamiCode.mockLaunch();
        expect(mockEmojiBlastsCancel).toHaveBeenCalled();
        expect(onActivate).toHaveBeenCalledTimes(1);
    });
    it("stops emoji blasts and konami listening when the returned function is called", () => {
        const stop = initializeKonamiEmojiBlast(vi.fn());
        mockKonamiCode.mockLaunch();
        stop();
        expect(mockEmojiBlastsCancel).toHaveBeenCalled();
        expect(mockKonamiCode.disable).toHaveBeenCalled();
    });
});
//# sourceMappingURL=initializeKonamiEmojiBlast.test.js.map