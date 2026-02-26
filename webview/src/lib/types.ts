export interface Np3OpenRequest {
	filePath: string;
}

export interface Np3Error {
	message: string;
	code: string;
	rawData?: string;
}

export interface ParameterDefinition {
	key: string;
	label: string;
	type: "continuous" | "discrete";
	min: number;
	max: number;
	step: number;
	defaultValue: number;
	group: string;
	unit?: string;
	options?: { label: string; value: number }[];
}

export interface ParameterValue {
	key: string;
	value: number;
}

export interface Np3OpenResponse {
	hash: string;
	recipe: Record<string, unknown>; // We will define a more precise type if needed later
	parameters: ParameterDefinition[];
}

export interface IpcPatchPayload {
	field: string;
	value: number;
}

export interface IpcResetPayload {
	field?: string;
}

// Global interface for all IPC messages
export interface IpcMessage {
	type: string;
	payload: Record<string, unknown> | Np3Error | Np3OpenResponse | IpcPatchPayload | IpcResetPayload;
}
