export function getNested(obj: unknown, path: string): any {
	if (!obj || !path) return undefined;
	const parts = path.split(".");
	let curr: any = obj;
	for (const part of parts) {
		if (curr === null || curr === undefined) return undefined;
		curr = curr[part];
	}
	return curr;
}

export function setNested(
	obj: Record<string, any>,
	path: string,
	value: any,
): Record<string, any> {
	if (!obj || !path) return obj;
	const parts = path.split(".");
	const newObj = { ...obj };
	let curr = newObj;

	for (let i = 0; i < parts.length - 1; i++) {
		const part = parts[i];
		curr[part] = curr[part] ? { ...curr[part] } : {};
		curr = curr[part];
	}

	curr[parts[parts.length - 1]] = value;
	return newObj;
}
