<script>
	import { api, requireAuth } from '$lib/api'
	import { session } from '$app/stores'

	import Modal from '$lib/Modal.svelte'
	import Tag from '$lib/Tag.svelte'
	import IconButton from '$lib/IconButton.svelte'
	import ChecklistItem from '$lib/ChecklistItem.svelte'

	import EditIcon from '$lib/icons/edit.svelte'
	import DeleteIcon from '$lib/icons/trash.svelte'
	import CartIcon from '$lib/icons/shopping-cart.svelte'
	import BadIcon from '$lib/icons/thumbs-down.svelte'
	import OkIcon from '$lib/icons/archive.svelte'
	import GoodIcon from '$lib/icons/thumbs-up.svelte'
	import GreatIcon from '$lib/icons/award.svelte'

	export let mediaList
	export let media

	let notes = media.notes
	let showEditModal = false

	async function mediaDelete() {
		if (requireAuth($session)) return;
		await api('DELETE', fetch, $session, `media/${media.id}`)
		mediaList = mediaList.filter((v) => v.id != media.id)
	}

	async function mediaToggleComplete(e) {
		if (requireAuth($session)) return;
		media.completed = e.currentTarget.checked
		media.cart = false
		media.completed_date = media.completed ? new Date() : null
		await api('PUT', fetch, $session, `media/${media.id}`, media)
	}

	async function mediaCart() {
		if (requireAuth($session)) return;
		media.cart = media.cart ? false : true
		await api('PUT', fetch, $session, `media/${media.id}`, media)
	}

	async function mediaNotes() {
		if (requireAuth($session)) return;
		media.notes = notes
		await api('PUT', fetch, $session, `media/${media.id}`, media)
		showEditModal = false
	}

	function mediaEdit() {
		if (requireAuth($session)) return;
		showEditModal = true
	}

	function mediaCancel() {
		if (requireAuth($session)) return;
		if (notes == media.notes || confirm('Exit now? You will lose your changes.')) {
			notes = media.notes
			showEditModal = false
		}
	}

	async function setRating(rating) {
		media.rating = media.rating !== rating ? rating : null
		await api('PUT', fetch, $session, `media/${media.id}`, media)
	}
</script>

<ChecklistItem checked={media.completed} on:change={mediaToggleComplete}>
	{media.description}

	<slot slot="tags">

		{#if media.related_link}
			<a class="related_link" href={media.related_link}>&#x1f855;</a>
		{/if}
		{#if media.comments}
			<a class="comments" href={media.comments}>&#x1F4AC;</a>
		{/if}
		<Tag value={media.media_type} />
		{#if media.completed && media.completed_date}
			<Tag value={new Date(media.completed_date)} />
		{/if}
		{#if media.notes && media.notes.length > 0}
			<Tag value={'notes'} fg="var(--yellow-dark)" bg="var(--yellow-light)" />
		{/if}
		{#if media.rating === 0}
			<Tag value="bad" fg="var(--red-dark)" bg="var(--red-light)" />
		{/if}
		{#if media.rating === 1}
			<Tag value="ok" fg="var(--orange-dark)" bg="var(--orange-light)" />
		{/if}
		{#if media.rating === 2}
			<Tag value="good" fg="var(--green-dark)" bg="var(--green-light)" />
		{/if}
		{#if media.rating === 3}
			<Tag value="great" fg="var(--blue-dark)" bg="var(--blue-light)" />
		{/if}
	</slot>

	<slot slot="actions">
		<IconButton on:click={mediaEdit} description="Edit media status"><EditIcon /></IconButton>
		<IconButton on:click={mediaDelete} description="Delete media"><DeleteIcon /></IconButton>
		{#if !media.completed}
			<IconButton on:click={mediaCart} description="Cart media"><CartIcon /></IconButton>
		{/if}
	</slot>
</ChecklistItem>

{#if showEditModal}
	<Modal on:exit={mediaCancel}>
		<p class="subtitle is-flat">{media.description}</p>
		<div class="paragraph">
			<IconButton on:click={() => setRating(0)} active={media.rating === 0}><BadIcon /></IconButton>
			<IconButton on:click={() => setRating(1)} active={media.rating === 1}><OkIcon /></IconButton>
			<IconButton on:click={() => setRating(2)} active={media.rating === 2}><GoodIcon /></IconButton>
			<IconButton on:click={() => setRating(3)} active={media.rating === 3}><GreatIcon /></IconButton>
		</div>
		<textarea bind:value={notes} />
		<div class="button-group">
			<button class="button" on:click={mediaCancel}>Cancel</button>
			<button class="button is-black" on:click={mediaNotes}>Submit</button>
		</div>
	</Modal>
{/if}

<style>
	.media-display {
		display: inline;
	}

	.button-group {
		text-align: right;
		margin-top: 1em;
	}
</style>
