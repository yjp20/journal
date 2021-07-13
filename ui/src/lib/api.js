const base = import.meta.env.VITE_API_BASE;

export class NetworkError extends Error {}

export async function api(method, fetch, session, resource, data) {
	let res;
	try {
		res = await fetch(`${base}/${resource}`, {
			credentials: 'include',
			method: method,
			body: data !== undefined ? JSON.stringify(data) : undefined,
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/json',
				Cookie: session.isSSR && session.token ? `token=${session.token}` : undefined
			}
		});
		if (!res.ok) throw new Error()
	} catch {
		throw new NetworkError();
	}

	try {
		const v = await res.json();
		if (!res.ok) {
			throw new Error(v.error);
		}
		return v;
	} catch {}
}
