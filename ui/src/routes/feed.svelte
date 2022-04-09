<script context="module">
	import { api, requireAuth } from '$lib/api'
	import { startOfDay, startOfWeek, endOfWeek, isSameWeek, addWeeks } from 'date-fns'

	export async function load({ fetch, page, session }) {
		const dateString = page.query.get('date')
		const date = dateString ? new Date(dateString) : new Date()
		const feedItemPromise = api('GET', fetch, session, 'feed', { end: endOfWeek(date).toISOString() })
		const feedSourcesPromise = api('GET', fetch, session, 'feedsource')
		return {
			props: {
				date: date,
				items: await feedItemPromise,
				sources: await feedSourcesPromise
			}
		}
	}
</script>

<script>
	import { goto } from '$app/navigation'
	import { session } from '$app/stores'
	import FeedSource from '$lib/FeedSource.svelte'
	import IconButton from '$lib/IconButton.svelte'
	import AddIcon from '$lib/icons/plus.svelte'

	export let date
	export let items
	export let sources

	let showSources = false

	const dateFormat = new Intl.DateTimeFormat('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	})

	const today = new Date()
	const lastWeek = startOfDay(addWeeks(new Date(), -1))

	function selectOther() {
		const ans = prompt('enter date as format MM/DD/YYYY', '')
		if (ans == null) return
		try {
			const date = new Date(ans)
			goto(formatDateLink(date))
		} catch (e) {
			alert('Bad date format')
		}
	}

	function formatDateLink(date) {
		return `/feed?date=${encodeURIComponent(date.toISOString())}`
	}

	function getHostname(url) {
		const o = new URL(url)
		return o.hostname
	}

	async function addToFeed(id) {
		if (requireAuth($session)) return
		const index = items.findIndex(item => item.id === id)
		items[index].added = true
		await api('POST', fetch, $session, 'feedsource/add', { id })
	}
</script>

<h1 class="title">feed</h1>

<div class="feed-header">
	<div class="date-controls">
		<div class="date-display">
			{dateFormat.format(startOfWeek(date))} - {dateFormat.format(endOfWeek(date))}
		</div>
		<div>
			<a href="/feed" class="no-link button" class:is-blue={isSameWeek(today, date)}> this week </a>
			<a href={formatDateLink(lastWeek)} class="no-link button" class:is-blue={isSameWeek(lastWeek, date)}>
				last week
			</a>
			<button class="button" on:click={selectOther}>other</button>
		</div>
	</div>
	{#if $session.token}
		<button class="button is-black" style="margin-left: auto" on:click={() => (showSources = true)}> sources </button>
	{/if}
</div>

{#each items as item, index}
	<li class="feeditem" class:is-added={item.added}>
		<div class="feeditem-label aside">{index}.</div>
		<div class="feeditem-content">
			<a href={item.related_link} class="feeditem-display">
				{item.description}{' '}
				<span class="feeditem-meta">({getHostname(item.related_link)})</span>
			</a>
		</div>
		<div class="feeditem-actions">
			{#if !item.added}
				<IconButton on:click={() => addToFeed(item.id)}><AddIcon /></IconButton>
			{/if}
		</div>
	</li>
{/each}

{#if showSources}
	<FeedSource on:exit={() => (showSources = false)} feedSources={sources} />
{/if}

<style>
	.feed-header {
		display: flex;
		align-items: start;
		margin: 1rem 0;
	}

	.date-display {
		margin-bottom: 0.5rem;
		font-size: 1.25em;
	}

	.feeditem {
		position: relative;
		display: flex;
		margin: 0.25rem 0;
	}

	.feeditem.is-added {
		font-weight: bold;
	}

	.feeditem-label {
		position: absolute;
		left: -2.75em;
		font-size: 0.875rem;
		width: 2.5em;
		margin-top: 0.1em;
		color: var(--greydark);
		flex-shrink: 0;
		text-align: right;
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
