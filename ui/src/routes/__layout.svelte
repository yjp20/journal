<script context="module">
	import Cookies from 'js-cookie'

	async function hash(message) {
		const msgBuffer = new TextEncoder().encode(message)
		const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer)
		const hashArray = Array.from(new Uint8Array(hashBuffer))
		const hashHex = hashArray.map((b) => b.toString(16).padStart(2, '0')).join('')
		return hashHex
	}
</script>

<script>
	import { session } from '$app/stores'
	import { api, NetworkError } from '$lib/api'
	import { root } from '$lib/config'
	import '../style.scss'

	let password = ''
	let passwordError = ''
	let checkMessage = ''

	async function login(e) {
		e.preventDefault()
		const hashed = await hash(password)
		Cookies.set('token', hashed, { domain: root })
		try {
			await api('POST', fetch, $session, 'token')
			$session.token = hashed
			passwordError = ''
			checkMessage = ''
		} catch (e) {
			if (e instanceof NetworkError) passwordError = 'Network error'
			else passwordError = 'Invalid password'
			Cookies.remove('token')
		}
	}

	async function check() {
		try {
			await api('POST', fetch, $session, 'token')
			checkMessage = 'Check succeeded'
		} catch (e) {
			if (e instanceof NetworkError) checkMessage = 'Network error'
			checkMessage = 'Check failed'
			Cookies.remove('token')
			$session.token = undefined
		}
	}
</script>

<div class="layout">
	<div class="side">
		<h1><b>journal</b></h1>
		<nav class="sidenav">
			<li class="sidenav-item"><a href="/">dashboard</a></li>
			<li class="sidenav-item"><a href="/todos">todos</a></li>
			<li class="sidenav-item"><a href="/media">media</a></li>
			<li class="sidenav-item"><a href="/feed">feed</a></li>
			<li class="sidenav-item"><a href="/graph">graph</a></li>
			<li class="sidenav-item"><a href="/login">login</a></li>
		</nav>
	</div>
	<main class="container">
		<slot />
	</main>
</div>

<style>
	:global(body) {
		line-height: 1.5;
	}

	.layout {
		display: flex;
		min-height: 100vh;
		padding: 2rem 0;
		box-sizing: border-box;
	}

	.side {
		padding: 0 2rem;
		width: 10rem;
		flex-shrink: 0;
	}

	.sidenav {
		margin-top: 0.5rem;
		margin-bottom: 1rem;
	}

	.sidenav-item {
		list-style: none;
	}

	.sidenav-item a {
		color: inherit;
	}

	@media screen and (max-width: 800px) {
		.layout {
			display: block;
		}

		.side {
			display: flex;
			width: unset;
			margin-bottom: 3em;
		}

		.sidenav {
			display: block;
			flex: 1;
			margin: 0 0.3em;
		}

		.sidenav-item {
			display: inline-block;
			margin-left: 0.5em;
		}

		.sidenav a {
			display: inline;
		}

		:global(.aside) {
			display: none;
		}
	}
</style>
