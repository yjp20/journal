<script context="module">
	import { api } from '$lib/api';

	export async function load({ fetch, page, session }) {
		const dateString = page.query.get("date")
		const date = dateString ? new Date(dateString) : new Date()
		const feedItemPromise = api('GET', fetch, session, 'feed', { end: date.toISOString() })
		const feedSourcesPromise = api('GET', fetch, session, 'feedsource')
		return {
			props: {
				date: date,
				feedItems: await feedItemPromise,
				feedSources: await feedSourcesPromise
			}
		};
	}
</script>

<script>
	import { startOfDay, startOfWeek, endOfWeek, isSameWeek, addWeeks } from 'date-fns';
	import { goto } from '$app/navigation';
	import { session } from '$app/stores';
	import FeedSource from '$lib/FeedSource.svelte';
	import IconButton from '$lib/IconButton.svelte';
	import AddIcon from "$lib/icons/plus.svelte";

	export let date;
	export let feedItems;
	export let feedSources;

	let showSources = false

	const dateFormat = new Intl.DateTimeFormat('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});

	const today = new Date();
	const lastWeek = startOfDay(addWeeks(new Date(), -1));

	function selectOther() {
		const ans = prompt('enter date as format MM/DD/YYYY', '');
		const date = new Date(ans);
		goto(formatDateLink(date))
	}

	function formatDateLink(date) {
		return `/feed?date=${encodeURIComponent(date.toISOString())}`
	}

	function getHostname(url) {
		const o = new URL(url)
		return o.hostname
	}

	async function addToFeed(id) {
		await api("POST", fetch, $session, "feedsource/add", { id });
	}
</script>

<h1 class="title">feed</h1>

<div class="feed-header">
	<div class="date-controls">
		<div class="date-display">
			{dateFormat.format(startOfWeek(date))} - {dateFormat.format(endOfWeek(date))}
		</div>
		<div>
			<a
				href="/feed"
				class="no-link button"
				class:is-blue={isSameWeek(today, date)}>
				this week
			</a>
			<a
				href={formatDateLink(lastWeek)}
				class="no-link button"
				class:is-blue={isSameWeek(lastWeek, date)}>
				last week
			</a>
			<button class="button" on:click={selectOther}>other</button>
		</div>
	</div>
	<button class="button is-black" style="margin-left: auto" on:click={() => (showSources = true)}>
		sources
	</button>
</div>

{#await feedItems then items}
	{#each items as feedItem}
		<div class="feeditem">
			<div class="feeditem-content">
				<a href={feedItem.related_link} class="feeditem-display">
					{feedItem.description}{" "}
					<span class="feeditem-meta">({getHostname(feedItem.related_link)})</span>
				</a>
			</div>
			<div class="feeditem-actions">
				<IconButton on:click={() => addToFeed(feedItem.id)}><AddIcon /></IconButton>
			</div>
		</div>
	{/each}
{/await}

{#if showSources}
	<FeedSource on:exit={() => (showSources = false)} feedSources={feedSources} />
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
		display: block;
	}

	.feeditem-display {
		color: var(--black);
	}

	.feeditem:hover .feeditem-display {
		text-decoration: underline;
	}

	.feeditem-meta {
		font-size: 0.75em;
		color: var(--greydark);
	}

	.feeditem-actions {
		margin-left: auto;
		visibility: hidden;
	}

	.feeditem:hover .feeditem-actions {
		visibility: visible;
	}
</style>
