import { fireEvent, render } from "@testing-library/svelte";
import { describe, expect, it, vi } from "vitest";
import type { ParameterDefinition } from "$lib/types";
import ParameterSliderUnit from "./ParameterSliderUnit.svelte";

const mockDefinition: ParameterDefinition = {
	key: "sharpening",
	label: "Sharpening",
	type: "continuous",
	min: 0,
	max: 100,
	step: 1,
	defaultValue: 50,
	group: "Detail",
};

describe("ParameterSliderUnit", () => {
	it("renders correctly with given definition", () => {
		const onchange = vi.fn();
		const { getByLabelText, getByRole } = render(ParameterSliderUnit, {
			definition: mockDefinition,
			value: 60,
			originalValue: 50,
			onchange,
		});

		// Label should be rendered
		expect(getByLabelText("Sharpening")).toBeDefined();

		// Input should reflect value
		const numInput = getByRole("spinbutton") as HTMLInputElement;
		expect(numInput).toBeDefined();
		expect(numInput.value).toBe("60");
		expect((numInput as HTMLInputElement).type).toBe("number");

		// Range slider should reflect value
		const rangeInputs = document.querySelectorAll('input[type="range"]');
		expect(rangeInputs.length).toBe(1);
		expect((rangeInputs[0] as HTMLInputElement).value).toBe("60");
	});

	it("fires onchange only on pointerup/change, not input", async () => {
		const onchange = vi.fn();
		render(ParameterSliderUnit, {
			definition: mockDefinition,
			value: 60,
			originalValue: 50,
			onchange,
		});

		const rangeInput = document.querySelector('input[type="range"]') as HTMLInputElement;

		// Simulate dragging (input event)
		await fireEvent.input(rangeInput, { target: { value: "75" } });
		expect(onchange).not.toHaveBeenCalled();

		// Simulate mouse release (pointerup)
		await fireEvent.pointerUp(rangeInput, { target: { value: "75" } });
		expect(onchange).toHaveBeenCalledWith("sharpening", 75);
	});

	it("resets to defaultValue on double-click", async () => {
		const onchange = vi.fn();
		const { getByLabelText } = render(ParameterSliderUnit, {
			definition: mockDefinition,
			value: 60,
			originalValue: 50,
			onchange,
		});

		const label = getByLabelText("Sharpening");
		await fireEvent.dblClick(label);

		expect(onchange).toHaveBeenCalledWith("sharpening", mockDefinition.defaultValue);
	});

	it("clamps numeric input to [min, max]", async () => {
		const onchange = vi.fn();
		const { getByRole } = render(ParameterSliderUnit, {
			definition: mockDefinition,
			value: 60,
			originalValue: 50,
			onchange,
		});

		const numInput = getByRole("spinbutton") as HTMLInputElement;

		// Try setting beyond max
		await fireEvent.input(numInput, { target: { value: "150" } });
		await fireEvent.blur(numInput);
		expect(onchange).toHaveBeenCalledWith("sharpening", mockDefinition.max); // clamped to 100

		// Try setting below min
		await fireEvent.input(numInput, { target: { value: "-10" } });
		await fireEvent.blur(numInput);
		expect(onchange).toHaveBeenCalledWith("sharpening", mockDefinition.min); // clamped to 0
	});

	it("shows dirty indicator only when value !== originalValue", () => {
		const onchange = vi.fn();

		// Render dirty
		const { unmount, queryByLabelText } = render(ParameterSliderUnit, {
			definition: mockDefinition,
			value: 60,
			originalValue: 50,
			onchange,
		});

		expect(queryByLabelText("Modified")).not.toBeNull();
		unmount();

		// Render clean
		const res2 = render(ParameterSliderUnit, {
			definition: mockDefinition,
			value: 50,
			originalValue: 50,
			onchange,
		});

		expect(res2.queryByLabelText("Modified")).toBeNull();
	});

	it("supports Home/End keyboard shortcuts for min/max", async () => {
		const onchange = vi.fn();
		render(ParameterSliderUnit, {
			definition: mockDefinition,
			value: 60,
			originalValue: 50,
			onchange,
		});

		const rangeInput = document.querySelector('input[type="range"]') as HTMLInputElement;

		// Press Home → should jump to min
		await fireEvent.keyDown(rangeInput, { key: "Home" });
		expect(onchange).toHaveBeenCalledWith("sharpening", mockDefinition.min);

		onchange.mockClear();

		// Press End → should jump to max
		await fireEvent.keyDown(rangeInput, { key: "End" });
		expect(onchange).toHaveBeenCalledWith("sharpening", mockDefinition.max);
	});
});
