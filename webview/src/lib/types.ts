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

export interface ToneCurvePoint {
	input: number;
	output: number;
}

export interface ColorAdjustment {
	hue: number;
	saturation: number;
	luminance: number;
}

export interface UniversalRecipe {
	exposure?: number;
	contrast?: number;
	highlights?: number;
	shadows?: number;
	vibrance?: number;
	saturation?: number;
	pointCurve?: ToneCurvePoint[];
	[key: string]: number | string | boolean | ToneCurvePoint[] | undefined | null;
}

export const NP3_SCHEMA_VERSION = 1;

export interface Np3OpenResponse {
	hash: string;
	recipe: UniversalRecipe;
	parameterDefinitions: ParameterDefinition[];
}

export interface Np3SaveAsRequest {
	filePath: string;
}

export interface Np3SaveAsResponse {
	filePath: string;
}

export interface CopyPastePayload {
	version: number;
	recipe: Partial<UniversalRecipe>;
}

export interface IpcPatchPayload {
	field: string;
	value: number | string;
}

export interface IpcResetPayload {
	field?: string;
}

export type IpcPayload =
	| Np3OpenRequest
	| Np3OpenResponse
	| Np3SaveAsRequest
	| Np3SaveAsResponse
	| CopyPastePayload
	| IpcPatchPayload
	| IpcResetPayload
	| Np3Error
	| { filePath: string }
	| { text: string };

export interface IpcMessage {
	type: string;
	payload: IpcPayload;
}
