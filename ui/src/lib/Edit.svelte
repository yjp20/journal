<script context="module">
	export function getStripped(text, tags) {
		if (tags === undefined) return text;
		const arr = text.split('');

		for (let tag of Object.values(tags)) {
			if (tag.start !== undefined && tag.end !== undefined && tag.start < tag.end) {
				for (let i = tag.start; i < tag.end; i++) {
					arr[i] = undefined;
				}
			}
		}

		return arr
			.filter((v) => v !== undefined)
			.join('')
			.trim();
	}
</script>

<script>
	import IconButton from '$lib/IconButton.svelte';
	import Tag from './Tag.svelte';
	import saveIcon from '$lib/icons/save.svg';
	import cancelIcon from '$lib/icons/x.svg';

	export let text;
	export let placeholder = undefined;
	export let tags = undefined;
	export let cancel = undefined;
	export let fg = {};
	export let bg = {};
</script>

<form class="edit" on:submit>
	<input class="edit-text" {placeholder} bind:value={text} />
	<div class="edit-controls">
		<div class="edit-tags">
			{#if tags !== undefined}
				{#if Object.keys(tags).length > 0}
					{#each Object.entries(tags) as tag}
						<Tag fg={fg[tag[0]]} bg={bg[tag[0]]} value={tag[1].value} />
					{/each}
				{:else}
					<Tag fg="#777" bg="#ccc" value="No Tags" />
				{/if}
			{/if}
		</div>
		{#if cancel}
			<IconButton type="button" on:click={cancel} description={'Cancel'} icon={cancelIcon} />
		{/if}
		<IconButton description={'Create'} icon={saveIcon} />
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
