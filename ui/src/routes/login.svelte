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

{#if $session.token}
	<form class="login box is-vertical" class:is-loggedin={Boolean($session.token)} on:submit={check}>
		<p class="paragraph">Logged in!</p>
		<p class="paragraph">
			<button class="button">Check</button>
		</p>
		{#if checkMessage}
			<p class="label">{checkMessage}</p>
		{/if}
	</form>
{:else}
	<form class="login box is-vertical" class:is-loggedin={Boolean($session.token)} on:submit={login}>
		<div class="field">
			<label class="label" for="password">Password:</label>
			<input class="input" name="password" type="password" bind:value={password} />
			{#if passwordError}
				<p class="label">{passwordError}</p>
			{/if}
		</div>
		<div class="field">
			<button class="button">Login</button>
		</div>
	</form>
{/if}
