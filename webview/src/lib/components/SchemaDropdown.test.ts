import { fireEvent, render } from "@testing-library/svelte";
import { describe, expect, it, vi } from "vitest";
import type { ParameterDefinition } from "$lib/types";
import SchemaDropdown from "./SchemaDropdown.svelte";

const mockDefinition: ParameterDefinition = {
    key: "pictureControlBase",
    label: "Picture Control Base",
    type: "discrete",
    min: 0,
    max: 2,
    step: 1,
    defaultValue: 0,
    group: "Picture Control",
    options: [
        { label: "Standard", value: 0 },
        { label: "Neutral", value: 1 },
        { label: "Vivid", value: 2 },
    ],
};

describe("SchemaDropdown", () => {
    it("renders options correctly from definition.options", () => {
        const onchange = vi.fn();
        const { getByLabelText, getByText } = render(SchemaDropdown, {
            definition: mockDefinition,
            value: 0,
            originalValue: 0,
            onchange,
        });

        // Label should be rendered
        expect(getByLabelText("Picture Control Base")).toBeDefined();

        // Should render the selected option text
        expect(getByText("Standard")).toBeDefined();
    });

    it("protective error state triggers when an invalid value is passed", () => {
        const onchange = vi.fn();
        const { container } = render(SchemaDropdown, {
            definition: mockDefinition,
            value: 99, // Invalid value not in options
            originalValue: 0,
            onchange,
        });

        const triggerButton = container.querySelector(
            '[data-slot="select-trigger"]',
        ) as HTMLButtonElement;
        expect(triggerButton).toBeDefined();
        const isDisabled =
            triggerButton.hasAttribute("disabled") ||
            triggerButton.hasAttribute("data-disabled");
        expect(isDisabled).toBe(true);
        expect(triggerButton.getAttribute("title")).toBe("Invalid value: 99");
    });

    it("displays selected label matching current value", () => {
        const onchange = vi.fn();
        const { getByText } = render(SchemaDropdown, {
            definition: mockDefinition,
            value: 2,
            originalValue: 0,
            onchange,
        });

        // "Vivid" is the label for value=2
        expect(getByText("Vivid")).toBeDefined();
    });

    it("resets to defaultValue on double-click", async () => {
        const onchange = vi.fn();
        const { getByText } = render(SchemaDropdown, {
            definition: mockDefinition,
            value: 2,
            originalValue: 0,
            onchange,
        });

        const labelEl = getByText("Picture Control Base");
        await fireEvent.dblClick(labelEl);

        expect(onchange).toHaveBeenCalledWith(
            "pictureControlBase",
            mockDefinition.defaultValue,
        );
    });

    it("shows dirty indicator only when value !== originalValue", () => {
        const onchange = vi.fn();

        // Render dirty
        const { unmount, queryByLabelText } = render(SchemaDropdown, {
            definition: mockDefinition,
            value: 2,
            originalValue: 0,
            onchange,
        });

        expect(queryByLabelText("Modified")).not.toBeNull();
        unmount();

        // Render clean
        const res2 = render(SchemaDropdown, {
            definition: mockDefinition,
            value: 0,
            originalValue: 0,
            onchange,
        });

        expect(res2.queryByLabelText("Modified")).toBeNull();
    });
});
