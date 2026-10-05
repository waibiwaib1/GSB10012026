import { describe, expect, it, vi } from "vitest";
import { z } from "zod";

import { createBase } from "./createBase.js";

const context = {
	options: {},
	take: vi.fn(),
};

describe("createStratumTemplate", () => {
	describe("prepare", () => {
		it("is undefined when neither the Base nor the template definition defines prepare", () => {
			const base = createBase({
				options: { name: z.string() },
			});
			const preset = base.createPreset({
				about: { name: "Test" },
				blocks: [],
			});

			const template = base.createStratumTemplate({
				presets: [preset],
			});

			expect(template.prepare).toBeUndefined();
		});

		it("uses the Base's prepare when the template definition does not define prepare", () => {
			const base = createBase({
				options: { name: z.string() },
				prepare: () => ({ name: "from-base" }),
			});
			const preset = base.createPreset({
				about: { name: "Test" },
				blocks: [],
			});

			const template = base.createStratumTemplate({
				presets: [preset],
			});

			expect(template.prepare?.(context)).toEqual({ name: "from-base" });
		});

		it("uses the template definition's prepare when the Base does not define prepare", () => {
			const base = createBase({
				options: { name: z.string() },
			});
			const preset = base.createPreset({
				about: { name: "Test" },
				blocks: [],
			});

			const template = base.createStratumTemplate({
				prepare: () => ({ preset: "test" }),
				presets: [preset],
			});

			expect(template.prepare?.(context)).toEqual({ preset: "test" });
		});

		it("composes the Base's and template definition's prepares when both exist", () => {
			const base = createBase({
				options: { name: z.string() },
				prepare: () => ({ name: "from-base" }),
			});
			const preset = base.createPreset({
				about: { name: "Test" },
				blocks: [],
			});

			const template = base.createStratumTemplate({
				prepare: () => ({ preset: "test" }),
				presets: [preset],
			});

			expect(template.prepare?.(context)).toEqual({
				name: "from-base",
				preset: "test",
			});
		});
	});
});
