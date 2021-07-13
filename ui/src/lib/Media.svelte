<script>
	import { api } from '$lib/api';
	import { session } from '$app/stores';

	import Modal from '$lib/Modal.svelte';
	import Tag from '$lib/Tag.svelte';
	import IconButton from '$lib/IconButton.svelte';
	import Checkbox from '$lib/Checkbox.svelte';
	import editIcon from '$lib/icons/edit.svg';
	import deleteIcon from '$lib/icons/trash.svg';
	import cartIcon from '$lib/icons/shopping-cart.svg';

	export let mediaList;
	export let media;

	let notes = media.notes;
	let showEditModal = false;

	async function mediaDelete() {
		await api('DELETE', fetch, $session, `media/${media.id}`);
		mediaList = mediaList.filter((v) => v.id != media.id);
	}

	async function mediaToggleComplete(e) {
		media.completed = e.currentTarget.checked;
		media.cart = false;
		media.completed_date = media.completed ? new Date() : null;
		await api('PUT', fetch, $session, `media/${media.id}`, media);
	}

	async function mediaCart() {
		media.cart = media.cart ? false : true;
		await api('PUT', fetch, $session, `media/${media.id}`, media);
	}

	async function mediaNotes() {
		media.notes = notes;
		await api('PUT', fetch, $session, `media/${media.id}`, media);
	}

	function mediaEdit() {
		showEditModal = true;
	}

	function mediaCancel() {
		notes = media.notes;
		showEditModal = false;
	}
</script>

<div class="media" class:is-completed={media.completed}>
	<div class="media-toggle">
		<Checkbox
			description="Mark media as {media.completed ? 'not done' : 'done'}"
			checked={media.completed}
			on:change={mediaToggleComplete}
		/>
	</div>
	<div class="media-content">
		<p class="media-display">{media.description}</p>
		{#if media.related_link}
			<a class="related_l" href={media.related_link}>&#x1f855;</a>
		{/if}
		{#if media.completed && media.completed_date}
			<Tag value={new Date(media.completed_date)} />
		{/if}
		<Tag value={media.media_type} />
	</div>
	<div class="media-actions">
		<IconButton icon={editIcon} on:click={mediaEdit} description="Edit media status" />
		<IconButton icon={deleteIcon} on:click={mediaDelete} description="Delete media" />
		{#if !media.completed}
			<IconButton icon={cartIcon} on:click={mediaCart} description="Cart media" />
		{/if}
	</div>
</div>

{#if showEditModal}
	<Modal>
		<p class="subtitle">{media.description}</p>
		<textarea bind:value={notes} />
		<div class="button-group">
			<button class="button" on:click={mediaCancel}>Cancel</button>
			<button class="button" on:click={mediaNotes}>Submit</button>
		</div>
	</Modal>
{/if}

<style>
	.media {
		display: flex;
	}

	.media.is-completed {
		opacity: 0.5;
	}

	.media-toggle {
		margin-top: 0.1em;
	}

	.media-content {
		display: flex;
		flex-wrap: wrap;
	}

	.media-actions {
		margin-left: auto;
		opacity: 0;
		margin-top: -0.25em;
	}

	.media:hover .media-actions {
		opacity: 1;
	}

	.media-display {
		margin-left: 0.25em;
	}

	.media:hover .media-display {
		background: var(--blue-light);
	}

	.button-group {
		text-align: right;
		margin-top: 1em;
	}
</style>
