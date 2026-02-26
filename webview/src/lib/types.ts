export interface Np3OpenRequest {
    filePath: string;
}

export interface Np3Error {
    message: string;
    code: string;
}

export interface Np3OpenResponse {
    hash: string;
    recipe: Record<string, unknown>; // We will define a more precise type if needed later
}

// Global interface for all IPC messages
export interface IpcMessage {
    type: string;
    payload: Record<string, unknown> | Np3Error | Np3OpenResponse;
}
