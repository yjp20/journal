<script>
	export let description;
	export let type = undefined;
	export let active = false;
</script>

<button on:click class="icon-button" class:is-active={active} aria-label={description} {type}>
	<slot />
</button>

<style>
	.icon-button {
		align-self: center;
		font-size: 1em;
		border: 0;
		padding: 0;
		background: none;
		margin-left: 0.125em;
		padding: 0.25em;
		border-radius: 5px;
	}

	.icon-button:hover {
		background-color: #f3f3f3;
	}

	.icon-button > :global(svg) {
		width: 1em;
		height: 1em;
		display: block;
	}

	.icon-button.is-active {
		color: white;
		background-color: black;
	}
</style>
