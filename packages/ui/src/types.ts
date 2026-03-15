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
