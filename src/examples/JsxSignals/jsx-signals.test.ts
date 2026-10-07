import { describe, expect, it } from "vitest";
import { screen } from "@testing-library/dom";
import { render } from "rezact";
import { Page } from "./JsxSignals";
import { delay } from "src/lib/utils";

describe("JSX Signals suite", () => {
  it("renders JSX Signals Page", async () => {
    render(document.body, Page);
    await delay(100);
    const allHeaders = screen.getAllByText(/jsx signals/i);
    expect(allHeaders).toHaveLength(1);
  });

  it("switches elements (method 1)", async () => {
    const clickToggleButton = screen.getByRole("button", {
      name: /change 1/i,
    });
    expect(clickToggleButton).not.toBeNull();
    expect(screen.getByText("Test Signal Element")).not.toBeNull();
    await clickToggleButton.click();
    expect(screen.getByText("Changed Signal Element")).not.toBeNull();
  });

  it("switches elements (method 2)", async () => {
    const clickToggleButton = screen.getByRole("button", {
      name: /change 2/i,
    });
    expect(clickToggleButton).not.toBeNull();
    expect(screen.getByText("Elm 1")).not.toBeNull();
    await clickToggleButton.click();
    expect(screen.getByText("Elm 2")).not.toBeNull();
    await clickToggleButton.click();
    expect(screen.getByText("Elm 1")).not.toBeNull();
    await clickToggleButton.click();
    expect(screen.getByText("Elm 2")).not.toBeNull();
  });

  it("switches elements (method 3)", async () => {
    const clickToggleButton = screen.getByRole("button", {
      name: /change 3/i,
    });
    expect(clickToggleButton).not.toBeNull();
    expect(screen.getByText("Elm 3")).not.toBeNull();
    await clickToggleButton.click();
    expect(screen.getByText("Elm 4")).not.toBeNull();
    await clickToggleButton.click();
    expect(screen.getByText("Elm 3")).not.toBeNull();
    await clickToggleButton.click();
    expect(screen.getByText("Elm 4")).not.toBeNull();
  });

  it("swaps elements and fragments", async () => {
    const changeToFragment = screen.getByRole("button", {
      name: /change 4/i,
    });
    const changeToElement = screen.getByRole("button", {
      name: /change 5/i,
    });
    expect(screen.getByText("Not Changed")).not.toBeNull();

    await changeToFragment.click();
    expect(screen.getByText("Fragment Part 1")).not.toBeNull();
    expect(screen.getByText("Fragment Part 2")).not.toBeNull();
    expect(screen.queryByText("Not Changed")).toBeNull();

    await changeToElement.click();
    expect(screen.getByText("Changed Back To Element")).not.toBeNull();
    expect(screen.queryByText("Fragment Part 1")).toBeNull();
    expect(screen.queryByText("Fragment Part 2")).toBeNull();

    await changeToFragment.click();
    expect(screen.getByText("Fragment Part 1")).not.toBeNull();
    expect(screen.getByText("Fragment Part 2")).not.toBeNull();
    expect(screen.queryByText("Changed Back To Element")).toBeNull();
  });
});
