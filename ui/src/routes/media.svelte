<script context="module">
	import { api } from '$lib/api';
	import Media from '$lib/Media.svelte';

	export async function load({ fetch, session }) {
		const mediaListPromise = api('GET', fetch, session, 'media')
		return {
			props: {
				mediaList: await mediaListPromise
			}
		};
	}
</script>

<script>
	import { flip } from 'svelte/animate';
	import { session } from '$app/stores';
	import Edit, { getStripped } from '$lib/Edit.svelte';

	export let mediaList;
	let text = '';
	let linkCache = {};

	const linkRegex = /https?:\/\/[.\/\-a-zA-Z0-9\_]*/;
	const mediaTypes = [
		[/:game/, 'game'],
		[/:book/, 'book'],
		[/:movie/, 'movie'],
		[/:show/, 'show'],
		[/:manga/, 'manga'],
		[/:anime/, 'anime']
	];

	async function createMedia(e) {
		e.preventDefault();
		const tags = await tagsPromise;
		if (!text) return;
		let media = {
			description: (tags["name"] && tags["name"].value) || getStripped(text, tags),
			related_link: tags["link"] && tags["link"].value,
			media_type: tags["media_type"] && tags['media_type'].value
		};
		const newMedia = await api('POST', fetch, $session, `media`, media);
		mediaList.push(newMedia);
		mediaList = mediaList;
		text = '';
	}

	async function getLinkMetadata(link) {
		if (linkCache[link] !== undefined) return linkCache[link];
		return linkCache[link] =
			await api('POST', fetch, $session, 'media/link', { link })
	}

	async function getTags(text) {
		const newTags = {};
		newTags['media_type'] = {
			value: 'articles',
			start: 0,
			end: 0
		};

		const matches = text.match(linkRegex);
		if (matches) {
			newTags['link'] = {
				value: matches[0],
				start: matches.index,
				end: matches.index + matches[0].length
			};
			const res = await getLinkMetadata(matches[0]);
			newTags['name'] = { value: res.name };
			newTags['media_type'] = { value: res.media_type };
		}

		for (let [pattern, name] of mediaTypes) {
			const matches = text.match(pattern);
			if (matches) {
				newTags['media_type'] = {
					value: name,
					start: matches.index,
					end: matches.index + matches[0].length
				};
			}
		}

		return newTags;
	}

	$: tagsPromise = getTags(text);
	$: sorted = mediaList.sort((a, b) => {
		if (a.completed !== b.completed) return a.completed - b.completed;
		if (a.completed_date !== b.completed_date)
			return new Date(b.completed_date) - new Date(a.completed_date);
		if (a.cart !== b.cart) return b.cart - a.cart;
		if (a.due_date !== b.due_date) {
			const a_date = a.due_date !== null ? new Date(a.due_date) : Infinity;
			const b_date = b.due_date !== null ? new Date(b.due_date) : Infinity;
			return a_date - b_date;
		}
		return a.id - b.id;
	});
</script>

<h1 class="title">media</h1>
<div class="mediaadd">
	<Edit placeholder="Media description or url" on:submit={createMedia} bind:text {tagsPromise} />
</div>

{#if sorted.length > 0 && !sorted[0].cart}
	<p class="paragraph"><em>No media in cart</em></p>
{/if}

{#each sorted as media (media.id)}
	<div
		animate:flip={{ duration: 100 }}
		class="media-wrapper"
		class:is-completed={media.completed}
		class:is-cart={media.cart}
	>
		<Media bind:mediaList bind:media />
	</div>
{/each}

<style>
	.mediaadd {
		margin-bottom: 1.5em;
	}

	.media-wrapper.is-cart:not(.is-completed) + .media-wrapper:not(.is-cart):not(.is-completed) {
		margin-top: 1rem;
	}

	.media-wrapper:not(.is-completed) + .media-wrapper.is-completed {
		margin-top: 1rem;
	}

	.media-wrapper + .media-wrapper {
		margin-top: 0.25rem;
	}
</style>
