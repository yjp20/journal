<script>
	import { format, formatRelative } from 'date-fns'

	export let fg = undefined
	export let bg = '#eee'
	export let value = undefined
</script>

<div class="tag" title={value instanceof Date ? format(value, "P p") : undefined} style="color: {fg}; background-color: {bg}">
	{#if value instanceof Date}
		{formatRelative(value, new Date())}
	{:else}
		{value.length > 20 ? `${value.substr(0, 30)}...` : value}
	{/if}
</div>

<style>
	.tag {
		display: inline-block;
		font-size: 0.75rem;
		padding: 0.1rem 0.3rem;
		border-radius: 5px;
	}
</style>
