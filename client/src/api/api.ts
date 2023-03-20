export type FromFormat = 'markdown'
export type ToFormat = 'docx' | 'html' | 'odt' | 'pdf'

export interface BaseResponse {
    type: 'StatusResponse' | 'ErrorResponse' | 'UploadResponse';
}

export interface UploadResponse extends BaseResponse {
    type: 'UploadResponse';
    uuid: string;
}

export interface StatusResponse extends BaseResponse {
    type: 'StatusResponse';
    status: 'uploaded' | 'converted';
}

export interface ErrorResponse extends BaseResponse {
    type: 'ErrorResponse';
    name: string;
    description: string;
}

export async function uploadFile(file: File, from: FromFormat, to: ToFormat): Promise<UploadResponse|ErrorResponse> {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("from", from);
    formData.append("to", to);

    const resp = await fetch(process.env.REACT_APP_ENDPOINT + "/api/v1/upload", {
        method: "POST",
        body: formData
    });

    const json = await resp.json()
    if (resp.ok) {
        return {
            type: 'UploadResponse',
            uuid: json["uuid"]
        };
    }

    return {
        type: 'ErrorResponse',
        name: json["name"],
        description: json["description"]
    };
}

export async function status(uuid: string): Promise<StatusResponse|ErrorResponse> {
    const params = new URLSearchParams({
        uuid: uuid
    });
    const url = `${process.env.REACT_APP_ENDPOINT}/api/v1/status?${params}`;

    const resp = await fetch(url, {
        method: "GET"
    });

    const json = await resp.json();
    if (resp.ok) {
        return {
            type: 'StatusResponse',
            status: json["status"]
        };
    }

    return {
        type: 'ErrorResponse',
        name: json["name"],
        description: json["description"]
    };
}