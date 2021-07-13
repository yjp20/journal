<script>
	export let icon;
	export let description;
	export let type = undefined;
</script>

<button on:click class="icon-button" aria-label={description} {type}>
	<img src={icon} alt={description} />
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

	.icon-button > img {
		height: 1em;
		display: block;
	}
</style>
