import cookie from 'cookie'
import { v4 as uuid } from '@lukeed/uuid'

export const handle = async ({ request, resolve }) => {
	const cookies = cookie.parse(request.headers.cookie || '')
	request.locals.token = cookies.token
	const response = await resolve(request)
	return response
}

export function getSession(request) {
	return {
		token: request.locals.token,
		isSSR: true
	}
}
