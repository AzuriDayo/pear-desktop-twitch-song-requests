import { afterEach, describe, expect, it, vi } from "vite-plus/test";
import { cleanup, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Provider } from "react-redux";
import { MemoryRouter } from "react-router";
import store from "../app/store";
import { Settings } from "./Settings";

afterEach(() => {
	cleanup();
	vi.unstubAllGlobals();
});

function renderSettings(enabled?: boolean, saveStatus = 200) {
	const fetchMock = vi.fn(async (url: string, init?: RequestInit) => {
		if (init?.method === "PATCH") {
			return new Response(JSON.stringify({ error: "Could not save setting" }), {
				status: saveStatus,
			});
		}
		return new Response(
			JSON.stringify(
				url === "/api/v1/settings" ? { pear_desktop_api_payload_dump_enabled: enabled } : [],
			),
		);
	});
	vi.stubGlobal("fetch", fetchMock);
	render(
		<Provider store={store}>
			<MemoryRouter>
				<Settings />
			</MemoryRouter>
		</Provider>,
	);
	return fetchMock;
}

describe("API payload dump setting", () => {
	it("defaults off when the setting is missing", async () => {
		renderSettings();
		const checkbox = screen.getByRole("checkbox", {
			name: "Enable Pear Desktop API payload dumps",
		});
		await waitFor(() => expect(checkbox).toBeEnabled());
		expect(checkbox).not.toBeChecked();
	});

	it.each([false, true])("loads %s and saves the opposite value", async (enabled) => {
		const user = userEvent.setup();
		const fetchMock = renderSettings(enabled);
		const checkbox = screen.getByRole("checkbox", {
			name: "Enable Pear Desktop API payload dumps",
		});
		await waitFor(() => expect(checkbox).toBeEnabled());
		expect((checkbox as HTMLInputElement).checked).toBe(enabled);
		await user.click(checkbox);
		await user.click(screen.getByRole("button", { name: "Save Debug Settings" }));
		await screen.findByText("Debug settings saved successfully!");
		expect(fetchMock).toHaveBeenCalledWith("/api/v1/settings", {
			method: "PATCH",
			body: JSON.stringify({ pear_desktop_api_payload_dump_enabled: String(!enabled) }),
		});
	});

	it("shows save failures", async () => {
		const user = userEvent.setup();
		renderSettings(false, 500);
		const checkbox = screen.getByRole("checkbox", {
			name: "Enable Pear Desktop API payload dumps",
		});
		await waitFor(() => expect(checkbox).toBeEnabled());
		await user.click(checkbox);
		await user.click(screen.getByRole("button", { name: "Save Debug Settings" }));
		await screen.findByText("Save failed: Could not save setting");
	});
});
