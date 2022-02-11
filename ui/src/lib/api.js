const base = import.meta.env.VITE_API_BASE

export class NetworkError extends Error {}

export async function api(method, fetch, session, resource, data) {
	let res
	try {
		if (method == 'GET') {
			const params = new URLSearchParams(data)
			const queryString = data ? `?${params.toString()}` : ''
			res = await fetch(`${base}/${resource}${queryString}`, {
				credentials: 'include',
				method: method,
				headers: {
					Accept: 'application/json',
					Cookie: session.isSSR && session.token ? `token=${session.token}` : undefined
				}
			})
		} else {
			res = await fetch(`${base}/${resource}`, {
				credentials: 'include',
				method: method,
				body: data !== undefined ? JSON.stringify(data) : undefined,
				headers: {
					Accept: 'application/json',
					'Content-Type': 'application/json',
					Cookie: session.isSSR && session.token ? `token=${session.token}` : undefined
				}
			})
		}
	} catch (e) {
		console.error(e)
		throw new NetworkError()
	}

	if (!res.ok) {
		throw new Error((await res.json()).error)
	}

	if (method.toLowerCase() !== 'delete' && method.toLowerCase() !== 'put') {
		let v
		try {
			v = await res.json()
		} catch (e) {
			throw new Error(`couldn't parse response as json`)
		}
		return v
	}
}

export function requireAuth(session) {
	if (!session.token) {
		alert('This function requires a login!')
		return true
	}
	return false
}
