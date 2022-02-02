<script context="module">
	export function getStripped(text, tags) {
		if (tags === undefined) return text
		const chars = text.split('')

		for (let tag of Object.values(tags)) {
			if (tag.start !== undefined && tag.end !== undefined && tag.start < tag.end) {
				for (let i = tag.start; i < tag.end; i++) {
					chars[i] = undefined;
				}
			}
		}

		return chars
			.filter((v) => v !== undefined)
			.join('')
			.trim()
	}
</script>

<script>
	import { createEventDispatcher } from 'svelte'

	const dispatch = createEventDispatcher()

	import IconButton from '$lib/IconButton.svelte'
	import Tag from './Tag.svelte'
	import SaveIcon from '$lib/icons/save.svelte'
	import CancelIcon from '$lib/icons/x.svelte'

	export let text
	export let placeholder = undefined
	export let tagsPromise = undefined
	export let cancel = false
	export let fg = {}
	export let bg = {}
</script>

<form class="edit" on:submit>
	<input class="edit-text" {placeholder} bind:value={text} />
	<div class="edit-controls">
		<div class="edit-tags">
			{#if tagsPromise !== undefined}
				{#await tagsPromise}
					Loading
				{:then tags}
					{#if Object.keys(tags).length > 0}
						{#each Object.entries(tags) as tag}
							<Tag fg={fg[tag[0]]} bg={bg[tag[0]]} value={tag[1].display || tag[1].value} />
						{/each}
					{/if}
				{/await}
			{/if}
		</div>
		{#if cancel}
			<IconButton type="button" on:click={() => dispatch('cancel')} description={'Cancel'}><CancelIcon /></IconButton>
		{/if}
		<IconButton description={'Create'}><SaveIcon /></IconButton>
	</div>
</form>

<style>
	.edit {
		width: 100%;
	}

	.edit-controls {
		display: flex;
		margin-top: 0.5rem;
	}

	.edit-tags {
		flex: 1;
	}
</style>
