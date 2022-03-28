<script>
	import { format } from 'date-fns'
	const date = new Date()
	export let records = [
		{
			date: date.toString(),
			description: 'drew a drawing :D on the wall',
			events: [
				{
					type: 'finish',
					description: 'Laundry'
				},
				{
					type: 'finish',
					description: 'Laundry'
				},
				{
					type: 'finish',
					description: 'Laundry'
				}
			]
		}
	]
</script>

<h1 class="title">dashboard</h1>
{#each records as record}
	<div class="record">
		<div class="record-header">
			<div class="record-date">
				<span class="label">{format(new Date(record.date), "d LLL ''yy")}</span>
			</div>
			<div class="record-desc">{record.description}</div>
		</div>
		<ul class="record-events">
			{#each record.events as event}
				<li class="record-event">{event.description}</li>
			{/each}
		</ul>
	</div>
{/each}

<style>
	.record {
		display: block;
	}

	.record-header {
		display: flex;
	}

	.record-date {
		display: inline-block;
		width: 4rem;
	}

	.record-events {
		margin-left: 2rem;
		font-size: 0.875em;
	}

	.record-event {
		list-style: disc;
	}
</style>
