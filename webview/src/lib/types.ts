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
	lane?: "left" | "right";
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

export interface ColorGradingZone {
	hue: number;
	chroma: number;
	brightness: number;
}

export interface ColorGrading {
	highlights: ColorGradingZone;
	midtone: ColorGradingZone;
	shadows: ColorGradingZone;
	blending: number;
	balance: number;
}

export interface UniversalRecipe {
	name?: string;
	description?: string;
	sourceFormat?: string;
	exposure?: number;
	contrast?: number;
	highlights?: number;
	shadows?: number;
	whites?: number;
	blacks?: number;
	vibrance?: number;
	saturation?: number;
	texture?: number;
	clarity?: number;
	dehaze?: number;
	sharpness?: number;
	sharpnessRadius?: number;
	sharpnessDetail?: number;
	sharpnessMasking?: number;
	midRangeSharpening?: number;
	temperature?: number | null;
	tint?: number;
	grainAmount?: number;
	grainSize?: number;
	grainRoughness?: number;
	red?: ColorAdjustment;
	orange?: ColorAdjustment;
	yellow?: ColorAdjustment;
	green?: ColorAdjustment;
	aqua?: ColorAdjustment;
	blue?: ColorAdjustment;
	purple?: ColorAdjustment;
	magenta?: ColorAdjustment;
	colorGrading?: ColorGrading;
	pointCurve?: ToneCurvePoint[];
	[key: string]:
		| number
		| string
		| boolean
		| ToneCurvePoint[]
		| ColorAdjustment
		| ColorGrading
		| undefined
		| null;
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

export interface ImageDataPayload {
	data: string;
	filename: string;
}

export interface WasmUriPayload {
	uri: string;
}

export type IpcPayload =
	| Np3OpenRequest
	| Np3OpenResponse
	| Np3SaveAsRequest
	| Np3SaveAsResponse
	| CopyPastePayload
	| IpcPatchPayload
	| IpcResetPayload
	| ImageDataPayload
	| WasmUriPayload
	| Np3Error
	| { filePath: string }
	| { text: string };

export interface IpcMessage {
	type: string;
	payload: IpcPayload;
	requestId?: string;
}
