<script>
	export let fg = undefined;
	export let bg = "#eee";
	export let value = undefined;
</script>

<div class="tag" style="color: {fg}; background-color: {bg}">
	{#if value instanceof Date}
		{new Intl.DateTimeFormat('en-US').format(value)}
	{:else}
		{value.length > 20 ? `${value.substr(0,30)}...` : value}
	{/if}
</div>

<style>
	.tag {
		display: inline-block;
		font-size: 0.75rem;
		padding: 0.2rem 0.3rem;
		border-radius: 5px;
	}

	:global(*) + .tag {
		margin-left: 0.25rem;
	}
</style>
