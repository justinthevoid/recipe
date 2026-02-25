export interface Np3OpenRequest {
    filePath: string;
}

export interface Np3OpenResponse {
    hash: string;
    recipe: Record<string, unknown>;
}

export interface Np3Error {
    type: "error";
    payload: {
        message: string;
        code: string;
    };
}
