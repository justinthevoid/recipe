interface VsCodeApi {
	postMessage(message: unknown): void;
	getState(): unknown;
	setState(state: unknown): void;
}

declare function acquireVsCodeApi(): VsCodeApi;

/**
 * Typed singleton wrapper for the VS Code webview API.
 * `acquireVsCodeApi()` can only be called once per webview lifecycle.
 */
class VsCodeWrapper {
	private readonly api: VsCodeApi;

	constructor() {
		this.api = acquireVsCodeApi();
	}

	postMessage(message: unknown): void {
		this.api.postMessage(message);
	}

	getState<T>(): T | undefined {
		return this.api.getState() as T | undefined;
	}

	setState<T>(state: T): void {
		this.api.setState(state);
	}
}

/** Singleton instance — safe because acquireVsCodeApi() is called exactly once. */
export const vscode = new VsCodeWrapper();
