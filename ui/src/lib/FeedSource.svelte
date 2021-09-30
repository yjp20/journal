<script>
	'@hmr:keep-all'

	import IconButton from '$lib/IconButton.svelte';
	import Modal from '$lib/Modal.svelte';
	import { api } from '$lib/api';
	import { session } from '$app/stores';
	import DeleteIcon from '$lib/icons/trash.svelte';

	export let feedSources;

	let description = '';
	let url = '';

	async function subscribeSourceFeed(e) {
		e.preventDefault();
		const source = await api('POST', fetch, session, 'feedsource', { description, url });
		feedSources = [...feedSources, source];
	}

	async function unsubscribeSourceFeed(id) {
		await api('DELETE', fetch, session, `feedsource/${id}`);
		feedSources = feedSources.filter(e => e.id != id);
	}
</script>

<Modal on:exit>
	<h2 class="title">sources</h2>
	{#each feedSources as source}
		<div class="feedsource">
			<div class="feedsource-description">{source.description}</div>
			<div class="feedsource-url">{source.url}</div>
			<div class="feedsource-remove">
				<IconButton on:click={() => unsubscribeSourceFeed(source.id)}><DeleteIcon /></IconButton>
			</div>
		</div>
	{/each}
	<form class="newfeedsource box is-vertical" on:submit={subscribeSourceFeed}>
		<div class="field">
			<label class="label" for="description">source description</label>
			<input description="description" bind:value={description} />
		</div>
		<div class="field">
			<label class="label" for="url">source url</label>
			<input description="url" bind:value={url} />
		</div>
		<div class="field">
			<button class="button">add source</button>
		</div>
	</form>
</Modal>

<style>
	.feedsource {
		position: relative;
		display: flex;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.feedsource-url {
		color: var(--greydark);
		font-size: 0.75rem;
		padding-left: 0.5rem;
	}

	.feedsource-remove {
		position: absolute;
		opacity: 0;
		margin-left: auto;
	}

	.feedsource:hover .feedsource-description {
		margin: -1px;
		padding: 1px;
		background-color: var(--blue-light);
		border-radius: 5px;
	}

	.feedsource:hover .feedsource-remove {
		opacity: 1;
		right: 0;
	}

	.newfeedsource {
		margin-top: 1rem;
	}
</style>
