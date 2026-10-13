import server from "./server";

export const GetPublicAgents = async (page: number, pageSize: number, key: string) => {
    const res = await server.Get("/workflow/public", { params: { page, pageSize, key } });
    return res;
};

export const GetPopularAgents = async (pageSize: number) => {
    const res = await server.Get("/workflow/popular", { params: { pageSize } });
    return res;
};