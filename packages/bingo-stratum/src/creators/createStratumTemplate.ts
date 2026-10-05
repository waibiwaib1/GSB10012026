import {
	AnyShape,
	InferredObject,
	LazyOptionalOptions,
	TemplatePrepareContext,
} from "bingo";
import { z } from "zod";

import { produceStratumTemplate } from "../producers/produceStratumTemplate.js";
import { Base } from "../types/bases.js";
import {
	StratumTemplate,
	StratumTemplateDefinition,
	StratumTemplateOptions,
	ZodPresetNameLiterals,
} from "../types/templates.js";
import { slugifyPresetName } from "../utils.ts/slugifyPresetName.js";

export function createStratumTemplate<OptionsShape extends AnyShape>(
	base: Base<OptionsShape>,
	templateDefinition: StratumTemplateDefinition<OptionsShape>,
): StratumTemplate<OptionsShape> {
	type Options = InferredObject<OptionsShape> & StratumTemplateOptions;

	const prepare =
		base.prepare || templateDefinition.prepare
			? (context: TemplatePrepareContext<Partial<Options>>) =>
					({
						...base.prepare?.(context),
						...templateDefinition.prepare?.(context),
					}) as LazyOptionalOptions<Partial<Options>>
			: undefined;
	const presetOption = z
		.union(
			templateDefinition.presets.map((preset) =>
				z.literal(slugifyPresetName(preset.about.name)),
			) as ZodPresetNameLiterals,
		)
		.describe("starting set of tooling to use");
	const template: StratumTemplate<OptionsShape> = {
		...templateDefinition,
		base,
		options: {
			...base.options,
			preset: presetOption.default(
				slugifyPresetName(
					(templateDefinition.suggested ?? templateDefinition.presets[0]).about
						.name,
				),
			) as unknown as z.ZodUnion<ZodPresetNameLiterals>, // TODO: why don't the types allow a ZodDefault here?
		},
		prepare,
		produce(context) {
			return produceStratumTemplate(template, context);
		},
	};

	return template;
}
