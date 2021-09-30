<script context="module">
	import { api } from '$lib/api';

	export async function load({ fetch, session }) {
		return {
			props: {
				initialDate: new Date(),
				initialFeedItems: await api('GET', fetch, session, 'feed'),
				initialFeedSources: await api('GET', fetch, session, 'feedsource')
			}
		};
	}
</script>

<script>
	import { startOfWeek, endOfWeek, isSameWeek, addWeeks } from 'date-fns';
	import { session } from '$app/stores';
	import FeedSource from '$lib/FeedSource.svelte';
	import IconButton from '$lib/IconButton.svelte';
	import AddIcon from "$lib/icons/plus.svelte";

	export let initialDate;
	export let initialFeedItems;
	export let initialFeedSources;

	const dateFormat = new Intl.DateTimeFormat('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});

	let selectedDate = initialDate;
	let feedItems = initialFeedItems;
	let showSources = false;

	const today = new Date();
	const lastWeek = addWeeks(new Date(), -1);

	async function collect() {
		await api('POST', fetch, $session, 'feed/collect');
	}

	function selectOther() {
		const ans = prompt('enter date as format MM/DD/YYYY', '');
		selectedDate = new Date(ans);
	}
</script>

<h1 class="title">feed</h1>

<div class="feed-header">
	<div class="date-controls">
		<div class="date-display">
			{dateFormat.format(startOfWeek(selectedDate))} - {dateFormat.format(endOfWeek(selectedDate))}
		</div>
		<div>
			<button
				class="button"
				class:is-blue={isSameWeek(today, selectedDate)}
				on:click={() => (selectedDate = today)}>this week</button
			>
			<button
				class="button"
				class:is-blue={isSameWeek(lastWeek, selectedDate)}
				on:click={() => (selectedDate = lastWeek)}>last week</button
			>
			<button class="button" on:click={selectOther}>other</button>
		</div>
	</div>
	<button class="button is-black" style="margin-left: auto" on:click={() => (showSources = true)}>
		sources
	</button>
</div>

{#each feedItems as feedItem}
	<div class="feeditem">
		<div class="feeditem-content">
			<div class="feeditem-display">{feedItem.description}</div>
			<div class="feeditem-meta">
				{#if feedItem.post_date}{feedItem.post_date} - {/if}{feedItem.related_link}
			</div>
		</div>
		<div class="feeditem-actions">
			<IconButton><AddIcon /></IconButton>
		</div>
	</div>
{/each}

{#if showSources}
	<FeedSource on:exit={() => (showSources = false)} feedSources={initialFeedSources} />
{/if}

<style>
	.feed-header {
		display: flex;
		align-items: start;
		margin: 1rem 0;
	}

	.date-display {
		margin-bottom: 8px;
		font-size: 1.25em;
	}

	.feeditem {
		position: relative;
		display: flex;
		margin: 0.75rem 0;
	}

	.feeditem-content {
		display: flex;
		flex-direction: column;
		align-items: start;
	}

	.feeditem-meta {
		font-size: 0.75em;
		color: var(--greydark);
	}

	.feeditem-actions {
		margin-left: auto;
		visibility: hidden;
	}

	.feeditem:hover .feeditem-display {
		margin: -1px;
		padding: 1px;
		background-color: var(--blue-light);
		border-radius: 5px;
	}

	.feeditem:hover .feeditem-actions {
		visibility: visible;
	}
</style>
